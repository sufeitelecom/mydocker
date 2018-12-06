package container

import (
	"os/exec"
	"syscall"
	"os"
	log "github.com/sirupsen/logrus"
	"strings"
	"io/ioutil"
)


//Namespace isolation
func Newprocess(tty bool,volume string) (*exec.Cmd, *os.File) {
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

	mnturl := "/home/sufei/busybox/root/mnt/"
	rooturl := "/home/sufei/busybox/root/"
	CreateWorkSpace(rooturl,mnturl,volume)
	cmd.Dir = mnturl

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
