package container

import (
	"os/exec"
	"syscall"
	"os"
	"github.com/sirupsen/logrus"
)


//Namespace isolation
func Newprocess(tty bool,command string) *exec.Cmd {
	logrus.Infof("Namespace isolation!")

	args := []string{"init",command}

	// create new process on the new Namespace
	cmd := exec.Command("/proc/self/exe",args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |  syscall.CLONE_NEWNET |syscall.CLONE_NEWUSER,
	}

	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: 0,
		Gid: 0,
	}

	cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 1001, Size: 1}}
	cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 1001, Size: 1}}

	if tty {
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
	}

	return cmd
}

//Init container.
func Runcontainerinit(command string, args []string) error {
	logrus.Infof("Init container ! command is %s",command)

	defaultmountflag := syscall.MS_NOEXEC |syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc","/proc","proc",uintptr(defaultmountflag),"")

	argv := []string{command}

	//exchange the initprocess , PID of command is 1
	if err := syscall.Exec(command,argv,os.Environ());err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}