package main

import (
	"fmt"
	"github.com/sufeitelecom/mydocker/container"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
	"encoding/json"
	"text/tabwriter"
)

func Listcontainers()  {
	dirs := fmt.Sprintf(container.DefaultInfoLocation,"")
	dirs = dirs[:len(dirs)-1]
	
	files,err := ioutil.ReadDir(dirs)
	if err != nil {
		log.Errorf("Read dir %s error. %v",dirs,err)
		return
	}

	var containers []*container.ContainerInfo
	for _,file := range files{
		if file.Name() == "network" {
			continue
		}
		tmpinfo,err := getContainerInfo(file)
		if err != nil{
			log.Errorf("Get container info error %v",err)
			continue
		}
		containers = append(containers,tmpinfo)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func getContainerInfo(file os.FileInfo) (*container.ContainerInfo,error) {
	containername := file.Name()
	configdir := fmt.Sprintf(container.DefaultInfoLocation,containername)
	configdir = configdir + "/" + container.ConfigName
	context,err := ioutil.ReadFile(configdir)
	if err != nil{
		log.Errorf("Read file % error %v",configdir,err)
		return nil,err
	}
	var containerinfo container.ContainerInfo
	if err := json.Unmarshal(context,&containerinfo);err != nil{
		log.Errorf("Json unmarshal error %v",err)
		return nil,err
	}
	return &containerinfo,nil
}
