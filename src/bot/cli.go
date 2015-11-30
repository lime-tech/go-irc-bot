package bot

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"path"
)

var (
	version string
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
	app.Commands = commands

	return app.Run(args)
}