package main

import (
	"github.com/sufeitelecom/mydocker/container"
	"fmt"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

func commitContainer(containername string,imagename string)  {
	path := fmt.Sprintf(container.MntUrl,containername)
	path += "/"

	image_tar :=  container.RootUrl + "/" + imagename + ".tar"
	fmt.Printf("%s",image_tar)
	if _,err := exec.Command("tar","-czf",image_tar,"-C",path,".").CombinedOutput();err != nil{
		log.Errorf("tar folder %s error %v",path,err)
	}
}
