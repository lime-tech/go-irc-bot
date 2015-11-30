package client

import (
	"bytes"
	"errors"
	log "github.com/Sirupsen/logrus"
	irc "github.com/fluffle/goirc/client"
	state "github.com/fluffle/goirc/state"
	"github.com/flynn-archive/go-shlex"
	"go-irc-bot/src/bot"
	"go-irc-bot/src/config"
	"strings"
	"time"
)

// Tip on long latency networks(posible implementation, not actual):
// Measure time between server pings, multiply it by 2
// Now it is a time of sended messages buffer.
// If connection is dead you message will be dead
// until you don't have any retry system.
// If client reconnects and after successful connection
// it discoveres some messages in this buffer, it will
// send those messages.

type Client struct {
	Conn     *irc.Conn
	Channels []string
	Config   *irc.Config
}

type Message struct {
	Kind    int
	Content string
	Nick    string
	Channel string
}

type Command struct {
	Line   *irc.Line
	Config *irc.Config
	Shlex  string
}

const (
	MSG_KIND_NULL = iota
	MSG_KIND_CHAN = iota
	MSG_KIND_PRIV = iota
	MSG_TRIM_SET  = " \n"
	HLIGHT_SEP    = ": "
)

var (
	RESTRICTED_ARGS = []string{"-u", "--u", "-i", "--i"}
)

func (self *Client) runHooks(k string, hs map[string][]string, line *irc.Line) (string, error) {
	if hooks, ok := hs[k]; ok {
		res := ""
		for _, hook := range hooks {
			if reply, err := (&Command{
				Line:   line,
				Config: self.Config,
				Shlex:  hook,
			}).Execute(); err != nil {
				res = err.Error()
				return res, nil
			} else {
				if len(res) > 0 {
					res += "\n" + reply
				} else {
					res = reply
				}
			}
		}
		return res, nil
	}
	return "", nil
}

func gotHighlighted(nick string, msg string) (bool, int, int) {
	if nick == msg {
		return true, 0, len(msg)
	}

	for _, sep := range []string{":", ",", " "} {
		if strings.HasPrefix(msg, nick+sep) {
			return true, 0, len(nick + sep)
		}
	}

	suff := " " + nick
	if strings.HasSuffix(msg, suff) {
		return true, len(msg) - len(suff), len(msg)
	}

	for _, sep := range []string{":", ",", " ", "!", "?"} {
		for _, suff := range []string{"?", "!", " ", ","} {
			compNick := sep + nick + suff
			index := strings.Index(msg, compNick)
			if index > 0 {
				return true, index, index + len(compNick)
			}
		}
	}

	return false, 0, 0
}

func (self *Message) String() string {
	switch self.Kind {
	case MSG_KIND_CHAN:
		prefix := ""
		if len(self.Nick) > 0 {
			prefix = self.Nick + HLIGHT_SEP
		}
		return prefix + self.Content
	case MSG_KIND_PRIV:
		return self.Content
	default:
		return ""
	}
}

func (self *Message) Send(conn *irc.Conn) error {
	to := ""
	switch self.Kind {
	case MSG_KIND_CHAN:
		to = self.Channel
		// Channels should begin with # and have one or more letter
		if len(self.Channel) < 2 {
			return errors.New("Invalid channel name")
		}
	case MSG_KIND_PRIV:
		to = self.Nick
	default:
		return errors.New("Unsupported message kind")
	}

	for _, rawStr := range strings.Split(self.String(), "\n") {
		reply := strings.Trim(rawStr, "\n")
		if len(reply) > 0 {
			conn.Privmsg(to, reply)
		}
	}

	// conn is not aware of any errors at this time :(
	return nil
}

func (c *Command) PredefinedParams() []string {
	return []string{
		"--user", c.Line.Nick,
		"--id", c.Line.Ident,
	}
}

func (c *Command) isRestrictedArg(arg string) bool {
	for _, p := range RESTRICTED_ARGS {
		if strings.HasPrefix(arg, p) {
			return true
		}
	}
	return false
}

func (c *Command) FilterArgs(args []string) []string {
	res := []string{}
	for _, arg := range args {
		if c.isRestrictedArg(arg) {
			continue
		}
		res = append(res, arg)
	}
	log.Debugf("Args after filtration %s", strings.Join(res, ":"))
	return res
}

func (c *Command) Execute() (string, error) {
	output := new(bytes.Buffer)
	args, err := shlex.Split(c.Shlex)
	if err != nil {
		return "", err
	}
	vargs := append(
		append([]string{c.Config.Me.Nick}, (c.PredefinedParams())...),
		c.FilterArgs(args)...,
	)

	if err := bot.Run(
		vargs,
		output,
	); err != nil {
		return "", err
	} else {
		return output.String(), nil
	}
}

func (self *Client) OnMsg(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) != 2 {
		log.Error("Something nasty happens, line.Args length is not 2, it is", len(line.Args))
		return
	}

	origin, rawMsg := line.Args[0], line.Args[1]
	msg := strings.Trim(rawMsg, MSG_TRIM_SET)
	log.Debug(origin, "->", msg)

	if len(origin) == 0 {
		log.Error("Something nasty happenes, got msg with zero length origin!")
		return
	}

	kind := MSG_KIND_NULL
	switch {
	case strings.HasPrefix(origin, "#"):
		kind = MSG_KIND_CHAN
	case len(origin) > 0:
		kind = MSG_KIND_PRIV
	default:
		kind = MSG_KIND_NULL
	}

	if kind != MSG_KIND_NULL {
		var resMsg string
		if highlighted, start, end := gotHighlighted(self.Config.Me.Nick, msg); highlighted {
			resMsg = msg[:start] + " " + msg[end:]
		} else {
			if kind == MSG_KIND_CHAN {
				// Not highlighted on channel -> reply is not required
				return
			}
			resMsg = msg
		}

		reply, err := (&Command{
			Line:   line,
			Config: self.Config,
			Shlex:  strings.Trim(resMsg, MSG_TRIM_SET),
		}).Execute()
		if err != nil {
			log.Error(err)
			reply = err.Error()
		}

		if len(reply) == 0 {
			return
		}
		err = (&Message{
			Kind:    kind,
			Content: reply,
			Nick:    line.Nick,
			Channel: origin,
		}).Send(conn)
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		log.Error("Something nasty happens, got msg with kind null!")
		return
	}
}

func (self *Client) OnJoin(conn *irc.Conn, line *irc.Line) {
	if line.Nick == self.Config.Me.Nick {
		return
	}
	if out, err := self.runHooks(irc.JOIN, bot.Hooks, line); err != nil {
		log.Error(err)
		return
	} else {
		if len(out) > 0 {
			(&Message{
				Kind:    MSG_KIND_PRIV,
				Nick:    line.Nick,
				Content: out,
				Channel: line.Args[0],
			}).Send(conn)
		}
	}
}

func New(c *config.Config) *Client {
	clientConfig := &irc.Config{
		Me: &state.Nick{
			Nick:  c.Nick,
			Ident: c.Nick,
			Name:  "",
		},
		Server:      c.Server,
		Pass:        c.ServerPassword,
		PingFreq:    3 * time.Minute,
		NewNick:     func(s string) string { return s + "_" },
		Recover:     (*irc.Conn).LogPanic,
		SplitLen:    350,
		Timeout:     60 * time.Second,
		Version:     "",
		QuitMessage: "",
		SSL:         false,
		Flood:       true,
	}

	conn := irc.Client(clientConfig)
	client := &Client{
		Conn:     conn,
		Channels: c.Channels,
		Config:   clientConfig,
	}

	conn.HandleFunc(irc.PRIVMSG, client.OnMsg)
	conn.HandleFunc(irc.JOIN, client.OnJoin)
	conn.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		for _, channel := range client.Channels {
			conn.Join(channel)
		}
	})

	return client
}
