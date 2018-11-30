package container

import (
	"os/exec"
	"syscall"
	"os"
	"github.com/sirupsen/logrus"
)

func Newprocess(tty bool,command string) *exec.Cmd {
	args := []string{"init",command}
	cmd := exec.Command("/proc/self/exe",args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |  syscall.CLONE_NEWNET,
	}

	if tty {
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func Runcontainerinit(command string, args []string) error {
	logrus.Infof("command is %s",command)

	defaultmountflag := syscall.MS_NOEXEC |syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc","/proc","proc",uintptr(defaultmountflag),"")

	argv := []string{command}
	if err := syscall.Exec(command,argv,os.Environ());err != nil {
		logrus.Errorf(err.Error())
	}
	return nil

}