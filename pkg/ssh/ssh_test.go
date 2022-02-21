package ssh

import (
	"testing"

	"github.com/fatih/color"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func TestConn(t *testing.T) {
	if utils.DetectCI() {
		return
	}

	conn, err := NewSSHConnection("local", "sudhanshu", "192.168.86.173", 22, "")
	utils.OK(t, err)

	wait := make(chan bool)

	copied, err := conn.Upload("test.data", "test.data")
	utils.OK(t, err)
	utils.Equals(t, int64(10), copied)

	go func(wait chan bool) {
		for {
			select {
			case stdout := <-conn.Stdout:
				if stdout.Completed {
					wait <- true
				}
				color.New(color.FgCyan).Println(stdout.Name, "|", stdout.Message)
			case stderr := <-conn.Stderr:
				color.New(color.FgRed).Println(stderr.Name, "|", stderr.Message)
			}
		}
	}(wait)

	status, err := conn.Run("ls -la")
	utils.OK(t, err)
	utils.Equals(t, -1, status)
	<-wait

	status, err = conn.Run("env")
	utils.OK(t, err)
	utils.Equals(t, -1, status)
	<-wait

	status, err = conn.Run("apt-get update")
	utils.OK(t, err)
	utils.Equals(t, true, status > 0)
	<-wait

	conn.Close()
}
