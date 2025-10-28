package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"mgo/internal/cmd"
	"mgo/internal/conf"
)

const (
	Version = "1.0"
	AppName = "mgo"
)

func init() {
	conf.App.Version = Version
	conf.App.Name = AppName
}

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.Usage = "a mgo service"
	app.Commands = []cli.Command{
		cmd.Web,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
