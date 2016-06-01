package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/server"
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

	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatal("Missing SLACK_TOKEN env variable")
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(5)

	server := server.New(args.host, args.port, token, deploy.NewStore())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("server is up and running at %s", server.Addr)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signals:
		log.Println("signal received, shutting down...")
		server.Shutdown()
	}
}
