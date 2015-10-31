package main

import (
	log "github.com/Sirupsen/logrus"
	"go-irc-bot/cli"
	"os"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		log.Fatalf("Got an error initial run %+v", err)
		defer os.Exit(1)
	}
}
