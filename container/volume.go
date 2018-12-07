package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"fmt"
)

func CreateWorkSpace(containerName string,volume string)  {
	CreateReadonlyLayer(RootUrl)
	CreateWriteLayer(containerName)
	CreateMountPiont(containerName)
	if volume != ""{
		var volumeurl []string
		volumeurl = strings.Split(volume,":")
		length := len(volumeurl)
		if length == 2 && volumeurl[0] != "" && volumeurl[1] != ""{
			MountVolume(containerName,volumeurl)
			log.Infof("%q",volumeurl)
		}else {
			log.Infof("Volume parameter input is not correct.")
		}
	}
}

func MountVolume(containerName string,volumeurl []string)   {
	parenturl :=  volumeurl[0]
	if err := os.Mkdir(parenturl,0777);err != nil{
		log.Infof("Mkdir dir %s error.%v",parenturl,err)
	}
    containerurl := fmt.Sprintf(MntUrl,containerName)
	sonurl := containerurl + "/" + volumeurl[1]
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
	busyboxurl := rootURL + "/busybox"
	busyboxtarurl := rootURL + "busybox.tar.gz"
	log.Infof("Readonly layer is %s",busyboxurl)
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

func CreateWriteLayer(containerName string)  {
	writeurl := fmt.Sprintf(WriteLayerUrl,containerName)
	if err := os.MkdirAll(writeurl,0777);err != nil{
		log.Errorf("Mkdir dir %s error. %v",writeurl,err)
	}
}

func CreateMountPiont(containerName string)  {
	mntURL := fmt.Sprintf(MntUrl,containerName)
	if err := os.MkdirAll(mntURL,0777);err != nil{
		log.Errorf("MKdir dir %s error. %v",mntURL,err)
	}

	tmpWriteLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	imagedir := RootUrl + "/busybox"
	dirs := "dirs=" + tmpWriteLayer + ":" + imagedir
	cmd := exec.Command("mount","-t","aufs","-o",dirs,"none",mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Errorf("%v",err)
	}else {
		log.Infof("mount root success！")
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

func DeleteWorkSpace(containerName string,volume string)  {
	DeleteMount(containerName,volume)
	DeleteWriteLayer(containerName)
}

func DeleteMount(containerName string,volume string)  {
	if volume != ""{
		var volumeurl []string
		volumeurl = strings.Split(volume,":")
		length := len(volumeurl)
		if length == 2 && volumeurl[0] != "" && volumeurl[1] != ""{
			sonurl := fmt.Sprintf(MntUrl,containerName)
			sonURL := sonurl + "/" + volumeurl[1]
			cmd := exec.Command("umount",sonURL)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if err := cmd.Run();err !=  nil{
				log.Errorf("Unmount dir % error.%v",sonurl,err)
			}
		}
	}
	mntURL := fmt.Sprintf(MntUrl,containerName)
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

func DeleteWriteLayer(containerName string)  {
	writeurl := fmt.Sprintf(WriteLayerUrl,containerName)
	if err := os.RemoveAll(writeurl); err != nil{
		log.Errorf("Remove dir %s error. %v",writeurl,err)
	}
}