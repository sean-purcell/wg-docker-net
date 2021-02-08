package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	driver, err := wg.NewDriver()
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	result := make(chan error, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	handler := network.NewHandler(driver)
	go func() {
		err := handler.ServeUnix(*socket, 0)
		result <- err
	}()
	log.Printf("Serving")

	select {
	case res := <-result:
		err = res
	case sig := <-stop:
		err = fmt.Errorf("Received signal: %v", sig)
	}

	delErr := driver.Delete()
	if delErr != nil {
		return fmt.Errorf("%v, failed to delete driver: %v", err, delErr)
	}
	return err
}
