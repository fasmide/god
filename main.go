package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fasmide/god/god"
	reaper "github.com/ramr/go-reaper"
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

	go reaper.Reap()

	daemon, err := god.Load(*cPath)
	if err != nil {
		fmt.Printf("god: unable to load configuration: %s\n", err)
		os.Exit(1)
	}

	err = daemon.Run()

	fmt.Printf("god: %s\n", err)
	os.Exit(0)
}
