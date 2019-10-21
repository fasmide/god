package main

import (
	"flag"
	"log"
	"os"

	"github.com/fasmide/god/god"
)

var (
	cPath = flag.String("path", "config.yml", "path to configuration")
)

func main() {
	flag.Parse()

	daemon, err := god.Load(*cPath)
	if err != nil {
		log.Fatalf("unable to load configuration: %s", err)
	}

	log.Printf("Well hello: %+v", daemon)

	err = daemon.Run()

	// error or no error - god is not ment to stop - ever
	log.Printf("God failed: %s", err)
	os.Exit(0)
}
