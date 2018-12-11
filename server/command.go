package main

import (
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/sufeitelecom/mydocker/container"
	"github.com/sufeitelecom/mydocker/cgroups/subsystems"
	"os"
	"github.com/sufeitelecom/mydocker/network"
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
	Usage:"Create a container with namespace and cgroups limit ie: mydocker run -ti [command]`",
	Flags:[]cli.Flag{
		cli.BoolFlag{
			Name:"ti",
			Usage:"enable tty",
		},
		cli.BoolFlag{
			Name:"d",
			Usage:"detach container",
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
			Name:"name",
			Usage:"container name",
		},
		cli.StringFlag{
			Name: "v",
			Usage: "volume",
		},
		cli.StringSliceFlag{
			Name:"e",
			Usage:"set environment",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name: "p",
			Usage: "port mapping",
		},
	},
	Action: func(c *cli.Context) error{
		if len(c.Args()) < 1{
			return fmt.Errorf("Missing container command parameter.")
		}

		var cmdArray []string
		for _, arg := range c.Args() {
			cmdArray = append(cmdArray, arg)
		}
		log.Infof("command is %v",cmdArray)

		tty := c.Bool("ti")
		detach := c.Bool("d")
		if tty && detach {
			return fmt.Errorf("ti and d parameter is not both provided.")
		}

		containerName := c.String("name")
		envSlice := c.StringSlice("e")
		netvork := c.String("net")
		portmapping := c.StringSlice("p")

		volume := c.String("v")
		resconf := &subsystems.ResourceConfig{
			MemoryLimit: c.String("m"),
			CpuShare: c.String("cpushare"),
			CpuSet: c.String("cpuset"),
		}

		Run(tty,cmdArray,resconf,volume,containerName,envSlice,netvork,portmapping)
		return nil
	},
}

// mydocker commit
var commitcommand = cli.Command{
	Name: "commit",
	Usage: "Tar a container into image.",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 2 {
			return fmt.Errorf("Container name and image name must be provided.")
		}
		container_name := c.Args().Get(0)
		image_name := c.Args().Get(1)
		commitContainer(container_name,image_name)
		return nil
	},
}

var listcommand = cli.Command{
	Name: "ps",
	Usage:"list all the container.",
	Action: func(c *cli.Context) error {
		Listcontainers()
		return nil
	},
}

var stopcommand = cli.Command{
	Name:"stop",
	Usage:"stop a container.",
	Action: func(c *cli.Context) error{
		if len(c.Args()) < 1{
			return fmt.Errorf("please tell the container name.")
		}
		containername := c.Args().Get(0)
		stopContainer(containername)
		return nil
	},
}

var removecommand  = cli.Command{
	Name:"rm",
	Usage:"remove unused container.",
	Action: func(c *cli.Context) error{
		if len(c.Args()) < 1{
			return fmt.Errorf("Missing container name.")
		}
		containername := c.Args().Get(0)
		removecontainer(containername)
		return nil
	},
}

var logcommand = cli.Command{
	Name:"logs",
	Usage:"print logs of container.",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1{
			return fmt.Errorf("Missing container name.")
		}
		containername := c.Args().Get(0)
		logcontainer(containername)
		return nil
	},
}

var execcommand = cli.Command{
	Name:"exec",
	Usage:"exec a command into container.",
	Action: func(c *cli.Context) error {
		if os.Getenv(ENV_EXEC_PID) != ""{
			log.Infof("pid callback %s",os.Getpid())
			return nil
		}
		if len(c.Args()) < 2{
			return fmt.Errorf("Missing container name or command")
		}
		containername := c.Args().Get(0)
		var commandArray []string
		for _,arg := range c.Args().Tail() {
			commandArray = append(commandArray,arg)
		}
		execContainer(containername,commandArray)
		return nil
	},
}

var networkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command {
		{
			Name: "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action:func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
				if err != nil {
					return fmt.Errorf("create network error: %+v", err)
				}
				return nil
			},
		},
		{
			Name: "list",
			Usage: "list container network",
			Action:func(context *cli.Context) error {
				network.Init()
				network.ListNetwork()
				return nil
			},
		},
		{
			Name: "remove",
			Usage: "remove container network",
			Action:func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}