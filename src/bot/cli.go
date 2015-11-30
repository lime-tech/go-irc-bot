package bot

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"path"
)

var (
	version   string
	rootFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "user",
			Usage:  "user that was run this application",
			EnvVar: "USER",
			Value:  "",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "id of user requesting the action",
			Value: "",
		},
		cli.StringFlag{
			Name:  "bucket",
			Usage: "storage namespace identifier",
			Value: "global",
		},
		cli.StringFlag{
			Name:  "channel",
			Usage: "source channel identifier in IRC style - #example",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "silent, s",
			Usage: "do not produce the noise",
		},
	}
)

func Run(args []string, output io.Writer) error {
	app := cli.NewApp()
	app.Name = path.Base(args[0])
	app.Usage = "Simple irc workflow bot"
	app.Version = fmt.Sprintf("%s", version)
	app.Authors = []cli.Author{
		cli.Author{Name: "bob", Email: ""},
	}

	if output == nil {
		output = os.Stdout
	}
	app.Writer = output
	app.Flags = rootFlags
	app.Commands = Commands

	return app.Run(args)
}
