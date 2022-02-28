package recipes

import (
	"time"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/pkg/cmd"
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
				logger.Infof("%s | %s", stdout.Name, stdout.Message)
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-cm.Stderr:
				logger.Infof("%s | %s", stderr.Name, stderr.Message)
				if stderr.Completed {
					wait <- true
				}
			}
		}
	}(wait)

	for _, command := range commands {
		logger.Announcef("%s | running <%s>", name, command)
		startTime := time.Now()
		err := cm.Run(command)
		if err != nil {
			logger.Errorf("%s | error in running command <%s>: %v", name, command, err)
			continue
		}
		<-wait
		logger.Successf("%s | command <%s> ended successfully in %s", name, command, time.Since(startTime))
	}

}
