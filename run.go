package mydoker

import (
	"mydoker/container"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run(tty bool,command string)  {
	parent := container.Newprocess(tty,command)
	if err := parent.Start();err != nil{
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
