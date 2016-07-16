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

var args struct {
	host string
	port int
}

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

	srv := server.New(args.host, args.port, slackToken, githubToken, stores.NewMemory())

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
