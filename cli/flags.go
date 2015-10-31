package cli

import (
	"github.com/codegangsta/cli"
)

var (
	rootFlags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: "log level(debug, info, warn, error, fatal, panic)",
		},
	}
	clientFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "path to client config file in TOML format",
			Value: "go-irc-bot.toml",
		},
	}
)
