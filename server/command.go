package main

import (
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/sufeitelecom/mydocker/container"
)

var initcommand = cli.Command{
	Name:"init",
	Usage:"Init container process",
	Action: func(c *cli.Context) error {
		log.Infof("Init come on")
		cmd := c.Args().Get(0)
		log.Infof("command is %s",cmd)
		err := container.Runcontainerinit(cmd,nil)
		return  err
	},
}

var runcommand  = cli.Command{
	Name:"run",
	Usage:"create a container",
	Flags:[]cli.Flag{
		cli.BoolFlag{
			Name:"ti",
			Usage:"enable tty",
		},
	},
	Action: func(c *cli.Context) error{
		if len(c.Args()) < 1{
			return fmt.Errorf("Missing container command")
		}
		cmd := c.Args().Get(0)
		tty := c.Bool("ti")
		Run(tty,cmd)
		return nil
	},
}

