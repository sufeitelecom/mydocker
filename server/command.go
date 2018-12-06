package main

import (
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/sufeitelecom/mydocker/container"
	"github.com/sufeitelecom/mydocker/cgroups/subsystems"
)

var initcommand = cli.Command{
	Name:"init",
	Usage:"Init container process",
	Action: func(c *cli.Context) error {
		log.Infof("Init come on")
		err := container.Runcontainerinit()
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name: "v",
			Usage: "volume",
		},
	},
	Action: func(c *cli.Context) error{
		if len(c.Args()) < 1{
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range c.Args() {
			cmdArray = append(cmdArray, arg)
		}
		log.Infof("command is %v",cmdArray)
		//cmdArray = cmdArray[1:]
		tty := c.Bool("ti")
		volume := c.String("v")
		resconf := &subsystems.ResourceConfig{
			MemoryLimit: c.String("m"),
			CpuShare: c.String("cpushare"),
			CpuSet: c.String("cpuset"),
		}

		Run(tty,cmdArray,resconf,volume)
		return nil
	},
}

