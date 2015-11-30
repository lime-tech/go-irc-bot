package ping

import (
	"github.com/codegangsta/cli"
)

func CliPing(c *cli.Context) {
	c.App.Writer.Write([]byte("Pong you!"))
}
