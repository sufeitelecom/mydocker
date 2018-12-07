package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/sufeitelecom/mydocker/container"
	"fmt"
	"os"
)

func removecontainer(containername string)  {
	containerinfo,err := GetContainerInfoFromName(containername)
	if err != nil{
		log.Errorf("Get container Info error %v",err)
		return
	}

	if containerinfo.Status != container.STOP {
		log.Errorf("Couldn't remove running container.")
		return
	}

	dirurl := fmt.Sprintf(container.DefaultInfoLocation,containername)
	if err := os.RemoveAll(dirurl);err != nil{
		log.Errorf("Remove file %s error %v",dirurl,err)
		return
	}
	container.DeleteWorkSpace(containername,containerinfo.Volume)
}
