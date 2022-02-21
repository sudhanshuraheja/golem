package recipes

import (
	"time"

	"github.com/sudhanshuraheja/golem/pkg/cmd"
	"github.com/sudhanshuraheja/golem/pkg/log"
)

type Cmd struct {
}

func (c *Cmd) Run(commands []string) {
	name := "local"
	cm := cmd.NewCmd(name)

	wait := make(chan bool)
	go func(wait chan bool) {
		for {
			select {
			case stdout := <-cm.Stdout:
				log.Infof("%s | %s", stdout.Name, stdout.Message)
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-cm.Stderr:
				log.Infof("%s | %s", stderr.Name, stderr.Message)
				if stderr.Completed {
					wait <- true
				}
			}
		}
	}(wait)

	for _, command := range commands {
		log.Announcef("%s | running <%s>", name, command)
		startTime := time.Now()
		err := cm.Run(command)
		if err != nil {
			log.Errorf("%s | error in running command <%s>: %v", name, command, err)
			continue
		}
		<-wait
		log.Successf("%s | command <%s> ended successfully in %s", name, command, time.Since(startTime))
	}

}
