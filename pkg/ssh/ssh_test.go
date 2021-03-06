package ssh

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

func TestSSH(t *testing.T) {
	if localutils.DetectCI() {
		return
	}

	log := logger.NewCLILogger(6, 8)

	conn, err := NewSSHConnection("thebatch", "sudhanshu", "192.168.86.173", 22, "")
	utils.Test().Nil(t, err)

	wait := make(chan bool)

	copied, err := conn.Upload("test.data", "test.data")
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, int64(10), copied)

	go func(log *logger.CLILogger, wait chan bool) {
		stdoutCh := conn.Stdout()
		stderrCh := conn.Stderr()
		for {
			select {
			case stdout := <-stdoutCh:
				if stdout.Completed {
					wait <- true
				}
				if stdout.Message != "" {
					log.Debug(stdout.Name).Msgf("%s", stdout.Message)
				}
			case stderr := <-stderrCh:
				if stderr.Completed {
					wait <- true
				}
				if stderr.Message != "" {
					log.Error(stderr.Name).Msgf("%s", stderr.Message)
				}
			}
		}
	}(log, wait)

	status, err := conn.Run("ls -la D*")
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, -1, status)
	<-wait

	status, err = conn.Run("env | grep PATH")
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, -1, status)
	<-wait

	conn.Close()
}
