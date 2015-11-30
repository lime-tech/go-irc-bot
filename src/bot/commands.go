package bot

import (
	"github.com/codegangsta/cli"
	"go-irc-bot/src/actions/ping"
	"go-irc-bot/src/actions/postpone"
)

var (
	commands = []cli.Command{
		{
			Name:      "ping",
			ShortName: "p",
			Usage:     "Test ping command",
			Action:    ping.CliPing,
		},
		{
			Name:      "postpone",
			ShortName: "pp",
			Usage:     "Postpone a message for user",
			Flags: []cli.Flag{
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
			},
			Action: postpone.CliPostpone,
		},
		{
			Name:      "postpone-remove",
			ShortName: "pr",
			Usage:     "removes postponed message",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "key, k",
					Usage: "message key",
					Value: -1,
				},
			},
			Action: postpone.CliPostponeRemove,
		},
		{
			Name:      "postpone-list",
			ShortName: "pl",
			Usage:     "view postponed message lists",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page, p",
					Usage: "page to show",
					Value: 0,
				},
			},
			Action: postpone.CliPostponeList,
		},
	}
)
