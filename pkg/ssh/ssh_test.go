package ssh

import (
	"os"
	"strconv"
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

func TestSSH(t *testing.T) {
	if localutils.DetectCI() {
		return
	}

	host := os.Getenv("GOLEM_SSH_TEST_HOST")
	if host == "" {
		t.Skip("integration SSH test skipped (set GOLEM_SSH_TEST_HOST to enable)")
	}
	user := os.Getenv("GOLEM_SSH_TEST_USER")
	if user == "" {
		user = os.Getenv("USER")
	}
	name := os.Getenv("GOLEM_SSH_TEST_NAME")
	if name == "" {
		name = host
	}
	port := 22
	if portStr := os.Getenv("GOLEM_SSH_TEST_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		} else {
			t.Fatalf("invalid GOLEM_SSH_TEST_PORT: %v", err)
		}
	}
	keyPath := os.Getenv("GOLEM_SSH_TEST_KEY_PATH")

	log := logger.NewCLILogger(6, 8)

	conn, err := NewSSHConnection(name, user, host, port, keyPath)
	if err != nil {
		t.Skipf("integration SSH test skipped due to connection error: %v", err)
	}

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
