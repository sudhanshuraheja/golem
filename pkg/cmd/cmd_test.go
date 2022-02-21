package cmd

import (
	"testing"

	"github.com/fatih/color"
	"github.com/sudhanshuraheja/golem/pkg/utils"
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
				color.New(color.FgCyan).Println(stdout.Name, "|", stdout.Message, "|", stdout.Completed)
			case stderr := <-c.Stderr:
				if stderr.Completed {
					wait <- true
				}
				color.New(color.FgRed).Println(stderr.Name, "||", stderr.Message, "|", stderr.Completed)
			}
		}
	}(wait)

	err := c.Run("ls -la")
	utils.OK(t, err)
	<-wait

	err = c.Run("ls -la ./../..")
	utils.OK(t, err)
	<-wait

	err = c.Run("asdfasdf")
	utils.Contains(t, "exit status", err.Error())
	<-wait
}
