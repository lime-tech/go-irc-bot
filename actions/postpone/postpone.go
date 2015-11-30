package postpone

import (
	log "github.com/Sirupsen/logrus"
)

type Message struct {
	Id   int    `orm:"auto"`
	From string `orm:"size(100)"`
	To   string `orm:"size(100)"`
	Data string `orm:"size(255)"`
}

func Do(m Message) error {
	log.Infof("Postpone %+v", m)
	return nil
}
