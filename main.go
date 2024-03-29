package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cego/god/god"
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

	if os.Getpid() == 1 {
		fmt.Println("god: use docker's --init functionality to enable reaping of zombie processes")
	}

	daemon, err := god.Load(*cPath)
	if err != nil {
		fmt.Printf("god: unable to load configuration: %s\n", err)
		os.Exit(1)
	}

	err = daemon.Run()
	if err != nil {
		fmt.Printf("god: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("god: shutdown complete\n")
	os.Exit(0)

}
