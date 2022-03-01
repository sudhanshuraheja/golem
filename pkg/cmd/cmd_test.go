package cmd

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
)

func TestCmd(t *testing.T) {
	log := logger.NewCLILogger(6, 8)
	c := NewCmd("test")

	wait := make(chan bool)
	go func(log *logger.CLILogger, wait chan bool) {
		for {
			select {
			case stdout := <-c.Stdout:
				if stdout.Completed {
					wait <- true
				}
				log.Debug(stdout.Name).Msgf("%s", stdout.Message)
			case stderr := <-c.Stderr:
				if stderr.Completed {
					wait <- true
				}
				log.Error(stderr.Name).Msgf("%s", stderr.Message)
			}
		}
	}(log, wait)

	err := c.Run("ls -la")
	utils.Test().Nil(t, err)
	<-wait

	err = c.Run("ls -la ./../..")
	utils.Test().Nil(t, err)
	<-wait

	err = c.Run("asdfasdf")
	utils.Test().Contains(t, err.Error(), "exit status")
	<-wait
}
