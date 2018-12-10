package main

import (
	"fmt"
	"github.com/sufeitelecom/mydocker/container"
	"io/ioutil"
	"encoding/json"
	"strings"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"os"
	_ "github.com/sufeitelecom/mydocker/nsenter"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func execContainer(containername string,commandArray []string)  {
	pid , err := GetContainerPidByName(containername)
	if err != nil{
		log.Errorf("Exec container getpid error %v",err)
		return
	}

	cmdstr := strings.Join(commandArray," ")
	log.Infof("container pid is %s",pid)
	log.Infof("command %s",cmdstr)

	cmd := exec.Command("/proc/self/exe","exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID,pid)
	os.Setenv(ENV_EXEC_CMD,cmdstr)

	containerEnvs := getEnvsByPid(pid)
	cmd.Env = append(os.Environ(),containerEnvs...)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containername, err)
	}


}

func GetContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + "/" + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s error %v", path, err)
		return nil
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}