package container

import (
	"fmt"
	"syscall"
	"os/exec"
	"os"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

//Init container.
func Runcontainerinit() error {
	log.Infof("Init container ! ")
	cmdArr := ReadUserCommand()
	if cmdArr == nil || len(cmdArr) == 0{
		return fmt.Errorf("get user command error!")
	}

	setMount()
	path ,err := exec.LookPath(cmdArr[0])
	if err != nil {
		log.Errorf("can not find the command %s",path)
		return err
	}

	log.Infof("find command :%s",path)

	//exchange the initprocess , PID of command is 1
	if err := syscall.Exec(path,cmdArr[0:],os.Environ());err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func pivotroot(root string) error  {
	if err := syscall.Mount(root,root,"bond",syscall.MS_BIND|syscall.MS_REC,"");err != nil{
		return fmt.Errorf("Mount rootfs error : %v",err)
	}

	pivotpath := filepath.Join(root,"/.pivot_root")
	if err := os.Mkdir(pivotpath,0777);err != nil{
		return err
	}


	if err := syscall.PivotRoot(root,pivotpath); err != nil{
		return err
	}

	if err := syscall.Chdir("/");err != nil{
		return fmt.Errorf("chdir / %v", err)
	}

	defaultmountflag := syscall.MS_NOEXEC |syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc","/proc","proc",uintptr(defaultmountflag),"");err != nil{
		return fmt.Errorf("Mount proc error %v",err)
	}

	pivotpath = filepath.Join("/",".pivot_root")
	if err := syscall.Unmount(pivotpath,syscall.MNT_DETACH);err != nil{
		return fmt.Errorf("Unmount fail:%v",err)
	}

	return os.Remove(pivotpath)
}

func setMount()  {
	pwd_path,err := os.Getwd()
	if err != nil{
		log.Errorf("Get current location error :%v",err)
		return
	}
	log.Infof("Current location is %s",pwd_path)


	if err := pivotroot(pwd_path);err != nil{
		log.Errorf("Pivot_root error:%v",err)
		return
	}

	if err := syscall.Mount("tmpfs","/dev","tmpfs",syscall.MS_NOSUID|syscall.MS_STRICTATIME,"mode=755");err != nil{
		log.Errorf("Mount dev error %v",err)
		return
	}
	return
}