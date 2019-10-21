package main

import (
	"flag"
	"log"
	"os"

	"github.com/fasmide/god/god"
)

var (
	cPath        = flag.String("path", "config.yml", "path to configuration")
	printExample = flag.Bool("example", false, "output example yml and exit")
)

func main() {
	flag.Parse()

	if *printExample {
		god.ExampleYml()
		os.Exit(0)
	}

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
