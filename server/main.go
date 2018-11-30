package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)


func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "mydocker is a simple container,just for funny!"

	app.Commands = []cli.Command{
		initcommand,
		runcommand,
	}

	app.Before = func(c *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args);err != nil{
		log.Fatal(err)
	}
}
