package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func CreateWorkSpace(rootURL string,mntURL string,volume string)  {
	CreateReadonlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPiont(rootURL,mntURL)
	if volume != ""{
		var volumeurl []string
		volumeurl = strings.Split(volume,":")
		length := len(volumeurl)
		if length == 2 && volumeurl[0] != "" && volumeurl[1] != ""{
			MountVolume(rootURL,mntURL,volumeurl)
			log.Infof("%q",volumeurl)
		}else {
			log.Infof("Volume parameter input is not correct.")
		}
	}
}

func MountVolume(rootURL string,mntURL string,volumeurl []string)   {
	parenturl := rootURL + volumeurl[0]
	if err := os.Mkdir(parenturl,0777);err != nil{
		log.Infof("Mkdir dir %s error.%v",parenturl,err)
	}

	sonurl := mntURL + volumeurl[1]
	if err := os.Mkdir(sonurl,0777);err != nil{
		log.Infof("Mkdir dir %s error.%v",sonurl,err)
	}

	dirs := "dirs=" + parenturl
	cmd := exec.Command("mount","-t","aufs","-o",dirs,"none",sonurl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Errorf("%v",err)
	}
}

func CreateReadonlyLayer(rootURL string)  {
	busyboxurl := rootURL + "busybox/"
	busyboxtarurl := rootURL + "busybox.tar.gz"
	exist ,err := ExistPath(busyboxurl)
	if err != nil{
		log.Infof("Fail to judge whether dir %s exist. %v",busyboxurl,err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxurl,0777);err != nil{
			log.Errorf("MKdir dir %s error. %v",busyboxurl,err)
		}
		if _,err := exec.Command("tar","-zxf",busyboxtarurl,"-C",busyboxurl).CombinedOutput();err != nil{
			log.Errorf("Untar dir %s error. %v",busyboxurl,err)
		}
	}
}

func CreateWriteLayer(rootURL string)  {
	writeurl := rootURL + "writelayer/"
	if err := os.Mkdir(writeurl,0777);err != nil{
		log.Errorf("Mkdir dir %s error. %v",writeurl,err)
	}
}

func CreateMountPiont(rootURL string,mntURL string)  {
	if err := os.Mkdir(mntURL,0777);err != nil{
		log.Errorf("MKdir dir %s error. %v",mntURL,err)
	}

	dirs := "dirs=" + rootURL + "writelayer:" + rootURL + "busybox"
	cmd := exec.Command("mount","-t","aufs","-o",dirs,"none",mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Errorf("%v",err)
	}
}

func ExistPath(path string) (bool,error)  {
	_,err := os.Stat(path)
	if err == nil{
		return true,nil
	}
	if os.IsNotExist(err) {
		return false,nil
	}
	return  false,err
}

func DeleteWorkSpace(rootURL string,mntURL string,volume string)  {
	DeleteMount(rootURL,mntURL,volume)
	DeleteWriteLayer(rootURL)
}

func DeleteMount(rootURL string,mntURL string,volume string)  {
	if volume != ""{
		var volumeurl []string
		volumeurl = strings.Split(volume,":")
		length := len(volumeurl)
		if length == 2 && volumeurl[0] != "" && volumeurl[1] != ""{
			sonurl := mntURL + volumeurl[1]
			cmd := exec.Command("umount",sonurl)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if err := cmd.Run();err !=  nil{
				log.Errorf("Unmount dir % error.%v",sonurl,err)
			}
		}
	}
	cmd := exec.Command("umount",mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Errorf("%v",err)
	}
	if err := os.RemoveAll(mntURL);err != nil{
		log.Errorf("Remove dir %s error. %v",mntURL,err)
	}
}

func DeleteWriteLayer(rootURL string)  {
	writeurl := rootURL + "writelayer/"
	if err := os.RemoveAll(writeurl); err != nil{
		log.Errorf("Remove dir %s error. %v",writeurl,err)
	}
}