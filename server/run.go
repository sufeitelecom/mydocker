package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/sufeitelecom/mydocker/container"
	"github.com/sufeitelecom/mydocker/cgroups"
	"github.com/sufeitelecom/mydocker/cgroups/subsystems"
	"time"
	"strings"
	"encoding/json"
	"fmt"
	"strconv"
)

func Run(tty bool,command []string,res *subsystems.ResourceConfig,volume string,containerName string,envSlice []string)  {

	containerID := container.RandStringBytes(10)
	if containerName == ""{
		containerName = containerID
	}

	parent,writepipe := container.Newprocess(tty,volume,containerName,envSlice)
	if parent == nil {
		log.Errorf("Create new process fail!")
		return
	}
	if err := parent.Start();err != nil{
		log.Error(err)
	}

	containerName, err := recordContainerInfo(parent.Process.Pid,containerID,containerName,command,volume)
	if err != nil{
		log.Errorf("Record container info error %v", err)
		return
	}

	cgroupmanager := cgroups.NewCgroupManager(containerID)
	defer cgroupmanager.Destory()
	cgroupmanager.Set(res)
	cgroupmanager.Apply(parent.Process.Pid)

	container.SendInitCommand(command,writepipe)

	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
		container.DeleteWorkSpace(containerName,volume)
	}

}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}

func recordContainerInfo(PID int,ID string,containerName string,commandArray []string,volume string) (string,error) {
	createtime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray," ")
	containerInfo := container.ContainerInfo{
		Name:containerName,
		Pid:strconv.Itoa(PID),
		Id:ID,
		Volume:volume,
		Status:container.RUNNING,
		Command:command,
		CreatedTime:createtime,
	}

	jsonbyte,err := json.Marshal(containerInfo)
	if err != nil{
		log.Errorf("Record container info error %v",err)
		return "",err
	}

	jsonstr := string(jsonbyte)
	dir := fmt.Sprintf(container.DefaultInfoLocation,containerName)
	if err := os.MkdirAll(dir,0622);err != nil{
		log.Errorf("Mkdir dir %s error. %v",dir,err)
		return  "",err
	}

	filename := dir + "/" + container.ConfigName
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", filename, err)
		return "", err
	}
	if _, err := file.WriteString(jsonstr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}
