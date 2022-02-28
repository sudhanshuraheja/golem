package cmd

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
)

func TestCmd(t *testing.T) {
	c := NewCmd("test")

	wait := make(chan bool)
	go func(wait chan bool) {
		for {
			select {
			case stdout := <-c.Stdout:
				if stdout.Completed {
					wait <- true
				}
				logger.Announcef("%s | %s", stdout.Name, stdout.Message)
			case stderr := <-c.Stderr:
				if stderr.Completed {
					wait <- true
				}
				logger.Errorf("%s | %s", stderr.Name, stderr.Message)
			}
		}
	}(wait)

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
