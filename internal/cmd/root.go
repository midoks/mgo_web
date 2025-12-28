package cmd

import (
	"fmt"
	"github.com/urfave/cli"

	"mgo/internal/conf"
	"mgo/internal/db"
	"mgo/internal/log"
)

var Root = cli.Command{
	Name:        "root",
	Usage:       "this command mod root password",
	Description: `mod root password`,
	Action:      runRoot,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "custom configuration file path"),
		stringFlag("password, p", "", "reset admin password"),
	},
}

// ./mgo_web root --password "admin"
func runRoot(c *cli.Context) error {
	conf.InitConf(c.String("config"))
	log.Init()
	db.InitDb()

	pwd := c.String("password")
	if pwd != "" {
		err := db.AdminUpdatePass(1, pwd)
		if err != nil {
			fmt.Println("update admin password fail!")
			return err
		}
		fmt.Println("update admin password success!")
	}
	return nil
}
