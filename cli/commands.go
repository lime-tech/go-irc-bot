package cli

import (
	"github.com/codegangsta/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "client",
			ShortName: "c",
			Usage:     "Run irc workflow client",
			Flags:     clientFlags,
			Action:    clientAction,
		},
	}
)
