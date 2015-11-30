package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"go-irc-bot/src/app"
	"go-irc-bot/src/config"
	"go-irc-bot/src/httpapi"
	"net/http"
	"os"
)

func clientAction(cx *cli.Context) {
	c, err := config.FromFile(cx.String("config"))
	if err != nil {
		log.Error(err)
		defer os.Exit(1)
		return
	}
	log.Debugf("Running client with config %+v", c)

	quit := make(chan bool)

	clients := app.RunClients(c)
	go func() {
		for n, client := range clients {
			router := httpapi.NewRouter(c.Http, client)
			http.Handle("/"+n, router)
		}
		err := http.ListenAndServe(c.Http.Addr, nil)
		if err != nil {
			log.Error(err)
			defer os.Exit(1)
			return
		}
	}()

	<-quit
}
