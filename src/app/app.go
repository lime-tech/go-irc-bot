package app

import (
	log "github.com/Sirupsen/logrus"
	irc "github.com/fluffle/goirc/client"
	"go-irc-bot/src/client"
	"go-irc-bot/src/config"
	"time"
)

func RunClients(cf *config.Config) map[string]*client.Client {
	res := map[string]*client.Client{}

	for n, c := range cf.Clients {
		cl := client.New(n, c)
		res[n] = cl
		conn := cl.Conn

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
					log.
						WithFields(log.Fields{"Profile": cl.Name, "Nick": cl}).
						Infof("Connected to %s", c.Server)
				}

				time.Sleep(2 * time.Second)
				<-reconnect
			}
		}()
	}
	return res
}
