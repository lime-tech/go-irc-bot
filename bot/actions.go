package bot

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func pingAction(c *cli.Context) {
	c.App.Writer.Write([]byte("Pong you!"))
}
