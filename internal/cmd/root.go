package cmd

import (
	"github.com/urfave/cli"

	// "mgo/internal/app"
	"mgo/internal/conf"
	"mgo/internal/db"
	"mgo/internal/log"
)

var Root = cli.Command{
	Name:        "root",
	Usage:       "this command start web services",
	Description: `start web services`,
	Action:      runRoot,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "custom configuration file path"),
	},
}

func runRoot(c *cli.Context) error {
	conf.InitConf(c.String("config"))
	log.Init()
	db.InitDb()

	return nil
}
