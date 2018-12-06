package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/sufeitelecom/mydocker/container"
	"github.com/sufeitelecom/mydocker/cgroups"
	"github.com/sufeitelecom/mydocker/cgroups/subsystems"
)

func Run(tty bool,command []string,res *subsystems.ResourceConfig,volume string)  {
	parent,writepipe := container.Newprocess(tty,volume)
	if parent == nil {
		log.Errorf("Create new process fail!")
		return
	}
	if err := parent.Start();err != nil{
		log.Error(err)
	}

	cgroupmanager := cgroups.NewCgroupManager("mydocker")
	defer cgroupmanager.Destory()

	cgroupmanager.Set(res)
	cgroupmanager.Apply(parent.Process.Pid)

	container.SendInitCommand(command,writepipe)

	parent.Wait()
	mnturl := "/home/sufei/busybox/root/mnt/"
	rooturl := "/home/sufei/busybox/root/"
	container.DeleteWorkSpace(rooturl,mnturl,volume)
	//syscall.Mount("proc","/proc","proc",syscall.MS_NODEV,"")
	os.Exit(-1)
}
