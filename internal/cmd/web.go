package cmd

import (
	"github.com/urfave/cli"

	"mgo/internal/app"
	"mgo/internal/conf"
	"mgo/internal/log"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "this command start web services",
	Description: `start web services`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "custom configuration file path"),
	},
}

func runWeb(c *cli.Context) error {
	_ = conf.Init("")
	log.Init()
	// db.InitDb()
	app.Run()
	return nil
}
