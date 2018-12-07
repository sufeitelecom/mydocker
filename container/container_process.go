package container

import (
	"os/exec"
	"syscall"
	"os"
	log "github.com/sirupsen/logrus"
	"strings"
	"io/ioutil"
	"math/rand"
	"time"
	"fmt"
)

var (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = "/home/sufei/busybox/root/mydocker/%s"
	ConfigName          string = "config.json"
	//ContainerLogFile    string = "container.log"
	RootUrl				string = "/home/sufei/busybox/root"
	MntUrl				string = "/home/sufei/busybox/root/mnt/%s"
	WriteLayerUrl 		string = "/home/sufei/busybox/root/writeLayer/%s"
)

type ContainerInfo struct {
	Pid         string `json:"pid"`        //容器的init进程在宿主机上的 PID
	Id          string `json:"id"`         //容器Id
	Name        string `json:"name"`       //容器名
	Command     string `json:"command"`    //容器内init运行命令
	CreatedTime string `json:"createTime"` //创建时间
	Status      string `json:"status"`     //容器的状态
	Volume      string `json:"volume"`     //容器的数据卷
} 

//Namespace isolation
func Newprocess(tty bool,volume string,containerName string) (*exec.Cmd, *os.File) {
	log.Infof("Namespace isolation!")

	readPipe,writePipe,err := NewPipe()
	if err != nil{
		log.Errorf("create Pipe fail")
		return nil,nil
	}

	// create new process on the new Namespace
	cmd := exec.Command("/proc/self/exe","init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |  syscall.CLONE_NEWNET |syscall.CLONE_NEWUSER,
	}

	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: 0,
		Gid: 0,
	}

	cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 0, Size: 1}}
	cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 0, Size: 1}}

	if tty {
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}


	CreateWorkSpace(containerName,volume)
	cmd.Dir = fmt.Sprintf(MntUrl,containerName)

	return cmd,writePipe
}
func NewPipe()(*os.File,*os.File,error)  {
	read, write,err := os.Pipe()
	if err != nil{
		return nil,nil,err
	}
	return read,write,err
}

func SendInitCommand(comarray []string,writepipe *os.File)  {
	command := strings.Join(comarray," ")
	log.Infof("command is %s",command)

	writepipe.WriteString(command)
	writepipe.Close()
}

func ReadUserCommand() []string  {
	pipe := os.NewFile(uintptr(3),"pipe")
	msg ,err := ioutil.ReadAll(pipe)

	if err != nil {
		log.Errorf("init read pipe fail: %v",err)
		return nil
	}

	msgStr := string(msg)
	return strings.Split(msgStr," ")
}

func RandStringBytes(n int) string  {
	letter := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte,n)
	for i := range b{
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}