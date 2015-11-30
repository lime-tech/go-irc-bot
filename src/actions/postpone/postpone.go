package postpone

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

const (
	LIMIT = 20
)

var (
	o           orm.Ormer
	CliCommands = []cli.Command{
		{
			Name:      "postpone",
			ShortName: "pp",
			Usage:     "Postpone a message for user",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "message text to postpone",
					Value: "",
				},
				cli.StringFlag{
					Name:  "to, t",
					Usage: "target user to postpone a message to",
					Value: "",
				},
			},
			Action: CliPostpone,
		},
		{
			Name:      "postpone-remove",
			ShortName: "pr",
			Usage:     "removes postponed message",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "key, k",
					Usage: "message key",
					Value: -1,
				},
			},
			Action: CliPostponeRemove,
		},
		{
			Name:      "postpone-list",
			ShortName: "pl",
			Usage:     "view postponed message lists",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page, p",
					Usage: "page to show",
					Value: 0,
				},
			},
			Action: CliPostponeList,
		},
	}
)

type Message struct {
	Id      int       `orm:"auto"`
	Bucket  string    `orm:"size(100);default(global)"`
	Channel string    `orm:"size(100)"`
	From    string    `orm:"size(100)"`
	To      string    `orm:"size(100)"`
	Data    string    `orm:"size(255)"`
	Date    time.Time `orm: "auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Message))
	if err := orm.RegisterDataBase("default", "sqlite3", "./postpone.db", 30); err != nil {
		log.Fatal(err)
	}
	if v, ok := os.LookupEnv("DB_RUN_SYNC"); ok && v == "yes" {
		if err := orm.RunSyncdb("default", true, true); err != nil {
			log.Fatal(err)
		}
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
		From:    c.GlobalString("user") + "@" + c.GlobalString("id"),
		Bucket:  c.GlobalString("bucket"),
		Channel: c.GlobalString("channel"),
		To:      to,
		Data:    message,
		Date:    time.Now(),
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
		Filter("bucket", c.GlobalString("bucket")).
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
		Filter("bucket", c.GlobalString("bucket")).
		Offset(c.Int("page") * LIMIT).
		Limit(LIMIT).
		All(ms); err != nil {
		log.Error(err)
		return
	} else {
		if !c.GlobalBool("silent") {
			c.App.Writer.Write([]byte(fmt.Sprintf("For user %s found %d postponed messages\n", user, num)))
		}

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
