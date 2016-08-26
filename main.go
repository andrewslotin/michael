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

	"github.com/andrewslotin/slack-deploy-command/auth"
	"github.com/andrewslotin/slack-deploy-command/bot"
	"github.com/andrewslotin/slack-deploy-command/dashboard"
	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/server"
	"github.com/andrewslotin/slack-deploy-command/slack"
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
	fmt.Println("Found a bug? Got an idea? Open an issue on https://github.com/andrewslotin/slack-deploy-command\nContributions are welcome!")
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
		slackBot.AddDeployEventHandler(bot.NewSlackTopicManager(api))
	} else {
		log.Printf("SLACK_WEBAPI_TOKEN env variable not set, channel topic notifications are disabled")
	}

	tokenSource := auth.RandomTokenSource{Src: rand.NewSource(time.Now().UnixNano())}
	authorizer := auth.NewOneTimeTokenAuthorizer(&tokenSource)

	slackBot.SetDashboardAuthorizer(authorizer)

	mux := http.NewServeMux()
	mux.Handle("/deploy", slackBot)
	mux.Handle("/", auth.TokenAuthMiddleware(deployDashboard, authorizer))

	srv := server.New(args.host, args.port)
	if err := srv.Start(mux); err != nil {
		log.Fatal(err)
	}

	log.Printf("slack-deploy-command server v%s is up and running at %s", version, srv.Addr)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signals:
		log.Println("signal received, shutting down...")
		srv.Shutdown()
	}
}
