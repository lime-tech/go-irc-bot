package bot

import (
	"github.com/codegangsta/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "ping",
			ShortName: "p",
			Usage:     "Test ping sction",
			Action:    pingAction,
		},
	}
)
