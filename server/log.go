package main

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/sufeitelecom/mydocker/container"
	"os"
	"io/ioutil"
)

func logcontainer(containername string)  {
	dirurl := fmt.Sprintf(container.DefaultInfoLocation,containername)
	filename := dirurl + "/" + container.ContainerLogFile

	file, err := os.Open(filename)
	if err != nil{
		log.Errorf("Open file %s error %v",filename,err)
		return
	}

	context , err := ioutil.ReadAll(file)
	if err != nil{
		log.Errorf("Read file %s error %v",filename,err)
		return
	}

	fmt.Fprint(os.Stdout,string(context))
}
