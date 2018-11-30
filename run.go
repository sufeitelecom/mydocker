package mydoker

import (

	log "github.com/sirupsen/logrus"
	"os"
	"github.com/sufeitelecom/mydocker/container"
)

func Run(tty bool,command string)  {
	parent := container.Newprocess(tty,command)
	if err := parent.Start();err != nil{
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
