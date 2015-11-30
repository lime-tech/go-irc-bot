package bot

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"go-irc-bot/src/actions/postpone"
)

func pingAction(c *cli.Context) {
	c.App.Writer.Write([]byte("Pong you!"))
}

func postponeAction(c *cli.Context) {
	to, message := c.String("to"), c.String("message")
	if len(to) == 0 || len(message) == 0 {
		cli.ShowCommandHelp(c, "postpone")
		return
	}
	if err := postpone.Do(postpone.Message{To: to, Data: message}); err != nil {
		log.Error(err)
		return
	}
	c.App.Writer.Write([]byte("Roger that!"))
}
