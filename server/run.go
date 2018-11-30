package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/sufeitelecom/mydocker/container"
)

func Run(tty bool,command string)  {
	parent := container.Newprocess(tty,command)
	if err := parent.Start();err != nil{
		log.Error(err)
	}
	parent.Wait()
	//syscall.Mount("proc","/proc","proc",syscall.MS_NODEV,"")
	os.Exit(-1)
}
