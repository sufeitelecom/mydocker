package main

import (
	"github.com/sufeitelecom/mydocker/container"
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"strconv"
	"syscall"
)

func stopContainer(containername string)  {
	coninfo ,err := GetContainerInfoFromName(containername)
	if err != nil{
		log.Errorf("Get container Info error %v",err)
		return
	}

	pidint,err := strconv.Atoi(coninfo.Pid)
	if err != nil{
		log.Errorf("Conver pid from string to int error %v",err)
		return
	}

	if err := syscall.Kill(pidint,syscall.SIGTERM);err != nil{
		log.Errorf("stop container %s error %v",containername,err)
		return
	}

	coninfo.Status = container.STOP
	coninfo.Pid =""

	newbyte,err := json.Marshal(coninfo)
	if err != nil{
		log.Errorf("Json marshal error %v",err)
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation,containername)
	config := dirURL + "/" + container.ConfigName
	if err := ioutil.WriteFile(config,newbyte,0622);err != nil{
		log.Errorf("Write file %s error %v",config,err)
		return
	}
	return
}

func GetContainerInfoFromName(containername string) (*container.ContainerInfo,error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation,containername)
	config := dirURL + "/" + container.ConfigName
	context ,err := ioutil.ReadFile(config)
	if err != nil{
		log.Errorf("Read file %s error %v",config,err)
	}

	var containerinfo container.ContainerInfo
	if err := json.Unmarshal(context,&containerinfo);err != nil{
		log.Errorf("Json unmarshal error %v",err)
		return nil,err
	}
	return &containerinfo,nil
}