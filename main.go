package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/docker/go-plugins-helpers/network"

	"github.com/iburinoc/wg-docker-net/wg"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var socket = flag.String("socket", "wg", "where to create the unix socket")
	flag.Parse()

	log.Printf("Creating socket at %s\n", *socket)

	stop := make(chan struct{})

	var driver = wg.NewDriver()
	var handler = network.NewHandler(driver)
	err := handler.ServeUnix(*socket, 0)
	if err != nil {
		return err
	}

	<-stop

	return nil
}
