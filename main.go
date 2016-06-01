package main

import (
	"flag"
	"fmt"
	"os"
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
}
