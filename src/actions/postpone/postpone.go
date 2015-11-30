package postpone

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const (
	LIMIT = 20
)

var (
	o orm.Ormer
)

type Message struct {
	Id   int       `orm:"auto"`
	From string    `orm:"size(100)"`
	To   string    `orm:"size(100)"`
	Data string    `orm:"size(255)"`
	Date time.Time ``
}

func init() {
	orm.RegisterModel(new(Message))
	if err := orm.RegisterDataBase("default", "sqlite3", "./postpone.db", 30); err != nil {
		log.Fatal(err)
	}
	if err := orm.RunSyncdb("default", true, true); err != nil {
		log.Fatal(err)
	}

	o = orm.NewOrm()
}

func Do(m Message) error {
	if id, err := o.Insert(&m); err != nil {
		return err
	} else {
		log.Infof("Postponed message, id %d", id)
	}

	return nil
}

func CliPostpone(c *cli.Context) {
	to, message := c.String("to"), c.String("message")
	if len(to) == 0 || len(message) == 0 {
		cli.ShowCommandHelp(c, "postpone")
		return
	}
	m := Message{
		From: c.GlobalString("user") + "@" + c.GlobalString("id"),
		To:   to,
		Data: message,
		Date: time.Now(),
	}
	if err := Do(m); err != nil {
		log.Error(err)
		return
	}
	c.App.Writer.Write([]byte("Roger that!"))
}

func CliPostponeRemove(c *cli.Context) {
	if c.Int("key") < 0 {
		c.App.Writer.Write(
			[]byte("--id should be positive integer value"),
		)
		return
	}

	user := c.GlobalString("user")
	qs := o.QueryTable("message")
	if num, err := qs.
		Filter("id", c.Int("key")).
		Filter("to", user).
		Delete(); err != nil {
		log.Error(err)
		return
	} else {
		c.App.Writer.Write(
			[]byte(
				fmt.Sprintf("Removed %d messages", num),
			),
		)
	}
}

func CliPostponeList(c *cli.Context) {
	qs := o.QueryTable("message")
	ms := new([]Message)
	user := c.GlobalString("user")
	if num, err := qs.
		Filter("to", user).
		Offset(c.Int("page") * LIMIT).
		Limit(LIMIT).
		All(ms); err != nil {
		log.Error(err)
		return
	} else {
		c.App.Writer.Write([]byte(fmt.Sprintf("For user %s found %d records\n", user, num)))
		for _, v := range *ms {
			c.App.Writer.Write(
				[]byte(
					fmt.Sprintf(
						"[%d] %s: %s> %s\n",
						v.Id,
						v.Date.String(),
						v.From,
						v.Data,
					),
				),
			)
		}
	}
}
