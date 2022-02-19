package ssh

import (
	"fmt"
	"testing"

	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func TestConn(t *testing.T) {
	conn, err := NewSSHConnection("sudhanshu", "192.168.86.173", 22, "")
	utils.OK(t, err)

	go func() {
		for {
			select {
			case stdout := <-conn.Stdout:
				fmt.Println(stdout)
			case stderr := <-conn.Stderr:
				fmt.Println("-----stderr-----", stderr)
			}
		}
	}()

	status, err := conn.Run("ls -la")
	utils.OK(t, err)
	utils.Equals(t, -1, status)
	<-conn.Completed

	status, err = conn.Run("env")
	utils.OK(t, err)
	utils.Equals(t, -1, status)
	<-conn.Completed

	status, err = conn.Run("apt-get update")
	utils.OK(t, err)
	utils.Equals(t, true, status > 0)
	<-conn.Completed

	conn.Close()
}
