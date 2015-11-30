package bot

import (
	"github.com/codegangsta/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "ping",
			ShortName: "p",
			Usage:     "Test ping command",
			Action:    pingAction,
		},
		{
			Name:      "postpone",
			ShortName: "pp",
			Usage:     "Postpone a message for user",
			Flags:     postponeFlags,
			Action:    postponeAction,
		},
	}
)
