package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	irc "github.com/fluffle/goirc/client"
	"go-irc-bot/src/client"
	"go-irc-bot/src/config"
	"go-irc-bot/src/httpapi"
	"net/http"
	"os"
	"time"
)

func clientAction(c *cli.Context) {
	cfg, err := config.FromFile(c.String("config"))
	if err != nil {
		log.Error(err)
		defer os.Exit(1)
		return
	}

	log.Debugf("Running client with config %+v", cfg)

	cClient := client.New(cfg)
	conn := cClient.Conn

	quit := make(chan bool)
	//incomming := make(chan *client.Message)
	go func() {
		for {
			reconnect := make(chan bool)
			conn.HandleFunc(irc.DISCONNECTED, func(_ *irc.Conn, _ *irc.Line) {
				reconnect <- true
			})

			if err := conn.Connect(); err != nil {
				log.Error(err)
				go func() { reconnect <- true }()
			} else {
				log.Infof("Connected to %s", cfg.Server)
			}

			time.Sleep(2 * time.Second)
			<-reconnect
		}
		panic("Unreachable")
	}()

	go func() {
		router := httpapi.NewRouter(cfg.Http, cClient)
		http.Handle("/", router)
		err := http.ListenAndServe(cfg.Http.Addr, nil)
		if err != nil {
			log.Error(err)
			defer os.Exit(1)
			return
		}
	}()

	<-quit
}
