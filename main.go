package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrewslotin/slack-deploy-command/deploy/stores"
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
		host string
		port int
	}
)

func init() {
	flag.StringVar(&args.host, "h", DefaultHost, "Host or address to listen on")
	flag.IntVar(&args.port, "p", DefaultPort, "Port to listen on")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

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

	store, err := getDeployStore()
	if err != nil {
		log.Fatalf("failed to open deploy DB: %s", err)
	}

	srv := server.New(args.host, args.port, slackToken, githubToken, store)

	if slackWebAPIToken := os.Getenv("SLACK_WEBAPI_TOKEN"); slackWebAPIToken != "" {
		api := slack.NewWebAPI(slackWebAPIToken, nil)
		srv.AddDeployEventHandler(server.NewSlackTopicManager(api))
	} else {
		log.Printf("SLACK_WEBAPI_TOKEN env variable not set, channel topic notifications are disabled")
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("server is up and running at %s", srv.Addr)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signals:
		log.Println("signal received, shutting down...")
		srv.Shutdown()
	}
}

func getDeployStore() (stores.Store, error) {
	boltDBPath := os.Getenv("BOLTDB_PATH")
	if boltDBPath != "" {
		log.Printf("writing deploy history into a BoltDB in %s", boltDBPath)
		return stores.NewBoltDB(boltDBPath)
	}

	log.Println("BOLTDB_PATH env variable not set, keeping deploy history in memory")

	return stores.NewMemory(), nil
}
