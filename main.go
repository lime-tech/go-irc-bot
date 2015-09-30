package main

import (
	log "github.com/Sirupsen/logrus"
	"irc-workflow/cli"
	"os"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		log.Fatalf("Got an error initial run %+v", err)
		defer os.Exit(1)
	}
}
