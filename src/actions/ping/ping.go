package ping

import (
	"github.com/codegangsta/cli"
)

var CliCommand = cli.Command{
	Name:      "ping",
	ShortName: "p",
	Usage:     "Test ping command",
	Action:    CliPing,
}

func CliPing(c *cli.Context) {
	c.App.Writer.Write([]byte("Pong you!"))
}
