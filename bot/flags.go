package bot

import (
	"github.com/codegangsta/cli"
)

var (
	rootFlags     = []cli.Flag{}
	postponeFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "message, m",
			Usage: "message text to postpone",
			Value: "",
		},
		cli.StringFlag{
			Name:  "to, t",
			Usage: "target user to postpone a message to",
			Value: "",
		},
	}
)
