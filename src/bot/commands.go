package bot

import (
	"github.com/codegangsta/cli"
	"go-irc-bot/src/actions/ping"
	"go-irc-bot/src/actions/postpone"
)

var (
	Commands = append(
		[]cli.Command{
			ping.CliCommand,
		},
		postpone.CliCommands...,
	)
	Hooks = map[string][]string{
		"JOIN": {"--silent postpone-list"},
	}
)
