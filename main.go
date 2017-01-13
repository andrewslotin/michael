package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrewslotin/michael/auth"
	"github.com/andrewslotin/michael/bot"
	"github.com/andrewslotin/michael/dashboard"
	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/server"
	"github.com/andrewslotin/michael/slack"
)

const (
	DefaultHost = "0.0.0.0"
	DefaultPort = 8081
)

var (
	binPath        = os.Args[0]
	version        = "n/a"
	buildDate      = "n/a"
	buildRev       = "n/a"
	buildGoVersion = "n/a"
	builder        = "n/a"

	args struct {
		host         string
		port         int
		printVersion bool
	}
)

func init() {
	flag.BoolVar(&args.printVersion, "version", false, "Print version and exit")
	flag.StringVar(&args.host, "h", DefaultHost, "Host or address to listen on")
	flag.IntVar(&args.port, "p", DefaultPort, "Port to listen on")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", binPath)
		flag.PrintDefaults()
	}
}

func printVersion() {
	fmt.Printf("%s v%s (rev %s)\n", binPath, version, buildRev)
	fmt.Printf("Built with %s on %s by %s\n\n", buildGoVersion, buildDate, builder)
	fmt.Println("Found a bug? Got an idea? Open an issue on https://github.com/andrewslotin/michael\nContributions are welcome!")
	os.Exit(0)
}

func main() {
	flag.Parse()

	if args.printVersion {
		printVersion()
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Fatal("Missing SLACK_TOKEN env variable")
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(5)

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Printf("GITHUB_TOKEN env variable not set, only public PRs details will be displayed in deploy announcements")
	}

	var (
		slackBot        *bot.Bot
		deployDashboard *dashboard.Dashboard
	)
	if boltDBPath := os.Getenv("BOLTDB_PATH"); boltDBPath != "" {
		log.Printf("writing deploy history into a BoltDB in %s", boltDBPath)

		store, err := deploy.NewBoltDBStore(boltDBPath)
		if err != nil {
			log.Fatalf("failed to open deploy DB: %s", err)
		}

		deployDashboard = dashboard.New(store)
		slackBot = bot.New(slackToken, githubToken, store)
	} else {
		log.Println("BOLTDB_PATH env variable not set, keeping deploy history in memory")

		store := deploy.NewInMemoryStore()
		deployDashboard = dashboard.New(store)
		slackBot = bot.New(slackToken, githubToken, store)
	}

	if slackWebAPIToken := os.Getenv("SLACK_WEBAPI_TOKEN"); slackWebAPIToken != "" {
		api := slack.NewWebAPI(slackWebAPIToken, nil)
		// Update channel topic to reflect current deploy status
		slackBot.AddDeployEventHandler(bot.NewSlackTopicManager(api))
		// Send direct messages to users mentioned in deploy subject
		slackBot.AddDeployEventHandler(bot.NewSlackIMNotifier(api))
	} else {
		log.Printf("SLACK_WEBAPI_TOKEN env variable not set, channel topic notifications are disabled")
	}

	tokenSource := auth.RandomTokenSource{Src: rand.NewSource(time.Now().UnixNano())}
	authenticator := auth.NewOneTimeTokenAuthenticator(&tokenSource)

	authSecret := os.Getenv("HISTORY_AUTH_SECRET")
	if authSecret == "" {
		authSecret = tokenSource.Generate(128)
		log.Printf("HISTORY_AUTH_SECRET is not set, using randomly generated %q", authSecret)
	}

	slackBot.SetDashboardAuth(authenticator)

	mux := http.NewServeMux()
	mux.Handle("/deploy", slackBot)
	mux.Handle("/", auth.TokenAuthenticationMiddleware(auth.ChannelAuthorizerMiddleware(deployDashboard, []byte(authSecret)), authenticator, []byte(authSecret)))

	srv := server.New(args.host, args.port)
	if err := srv.Start(mux); err != nil {
		log.Fatal(err)
	}

	log.Printf("Michael Buffer v%s is listening on %s", version, srv.Addr)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signals:
		log.Println("signal received, shutting down...")
		srv.Shutdown()
	}
}
