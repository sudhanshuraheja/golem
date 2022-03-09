package execssh

import (
	"sync"
	"testing"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/pkg/ssh"
)

type SSHMock struct {
	mu     sync.Mutex
	stdout chan ssh.Out
	stderr chan ssh.Out
}

func (m *SSHMock) Stdout() chan ssh.Out {
	m.mu.Lock()
	if m.stdout == nil {
		m.stdout = make(chan ssh.Out, 10)
	}

	m.mu.Unlock()
	return m.stdout
}

func (m *SSHMock) Stderr() chan ssh.Out {
	m.stderr = make(chan ssh.Out, 10)
	return m.stderr
}

func (m *SSHMock) Run(cmd string) (int, error) {
	m.mu.Lock()
	if m.stdout == nil {
		m.stdout = make(chan ssh.Out, 10)
	}
	m.stdout <- ssh.Out{
		Name:      "mock",
		ID:        "mock",
		Command:   cmd,
		Message:   "completed",
		Completed: true,
	}
	m.mu.Unlock()
	return 0, nil
}

func (m *SSHMock) Upload(src, dest string) (int64, error) {
	return int64(0), nil
}

func (m *SSHMock) Close() {

}

func TestExecSSH(t *testing.T) {
	mSSH := SSH{
		conn:   &SSHMock{},
		log:    logger.NewCLILogger(6, 8),
		name:   "mock",
		output: []ssh.Out{},
	}

	cmds := commands.Commands{}
	cmds.Append(commands.NewCommand("ls -la"))
	mSSH.Run(cmds)

	source := "source"
	art := artifacts.Artifact{
		Source:      &source,
		Destination: "destination",
	}
	mSSH.Upload([]*artifacts.Artifact{&art})
}
