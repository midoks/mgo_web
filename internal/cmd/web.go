package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"mgo/internal/app"
	"mgo/internal/conf"
	"mgo/internal/db"
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
	fmt.Println("1. Before InitConf")
	conf.InitConf(c.String("config"))
	fmt.Println("2. After InitConf, InstallLock:", conf.Security.InstallLock)

	log.Init()
	fmt.Println("3. After log.Init")

	if conf.Security.InstallLock {
		db.InitDb()
		fmt.Println("4. After db.InitDb")
	}

	app.Run()
	return nil
}
