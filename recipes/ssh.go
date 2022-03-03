package recipes

import (
	"fmt"
	"time"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
	"github.com/sudhanshuraheja/golem/pkg/ssh"
)

type SSH struct {
	conn   *ssh.Connection
	log    *logger.CLILogger
	name   string
	output []ssh.Out
}

type SSHJob struct {
	Server config.Server
	Recipe *Recipe
}

func (ss *SSH) Connect(s *config.Server) error {
	var host string
	switch {
	case s.PublicIP != nil:
		host = *s.PublicIP
	case len(s.HostName) > 0:
		host = s.HostName[0]
	default:
		return fmt.Errorf("could not find a public ip or a hostname in config")
	}

	startTime := time.Now()
	conn, err := ssh.NewSSHConnection(s.Name, s.User, host, s.Port, "")
	if err != nil {
		return fmt.Errorf("could not ssh to host: %v", err)
	}
	ss.log.Info(s.Name).Msgf("connected via SSH %s", localutils.TimeInSecs(startTime))
	ss.conn = conn
	ss.name = s.Name
	return nil
}

func (ss *SSH) Run(commands []string) {
	ss.output = []ssh.Out{}
	wait := make(chan bool)
	go func(log *logger.CLILogger, wait chan bool) {
		for {
			select {
			case stdout := <-ss.conn.Stdout:
				ss.output = append(ss.output, stdout)
				if stdout.Message != "" {
					ss.log.Debug(stdout.Name).Msgf("%s", stdout.Message)
				}
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-ss.conn.Stderr:
				ss.output = append(ss.output, stderr)
				if stderr.Message != "" {
					ss.log.Error(stderr.Name).Msgf("%s", stderr.Message)
				}
				if stderr.Completed {
					wait <- true
				}
			}
		}
	}(ss.log, wait)

	for _, command := range commands {
		startTime := time.Now()
		ss.log.Highlight(ss.name).Msgf("$ %s", command)
		status, err := ss.conn.Run(command)
		if err != nil {
			ss.log.Error(ss.name).Msgf("error in running command <%s>: %v", command, err)
			continue
		}
		<-wait
		if status > 0 {
			ss.log.Error(ss.name).Msgf("command <%s> failed with status: %d", command, status)
			continue
		}
		ss.log.Success(ss.name).Msgf("$ %s %s", command, localutils.TimeInSecs(startTime))
	}
}

func (ss *SSH) Upload(artifacts []config.Artifact) {
	for _, artifact := range artifacts {
		startTime := time.Now()
		ss.log.Info(ss.name).Msgf(
			"%s %s %s %s:%s",
			logger.Cyan("uploading"),
			localutils.TinyString(*artifact.Source, tiny),
			logger.Cyan("to"),
			ss.name,
			localutils.TinyString(artifact.Destination, tiny),
		)
		copied, err := ss.conn.Upload(*artifact.Source, artifact.Destination)
		if err != nil {
			ss.log.Error(ss.name).Msgf("failed to upload local:<%s> to remote:<%s>: %v", *artifact.Source, artifact.Destination, err)
			continue
		}
		ss.log.Success(ss.name).Msgf("uploaded %s to %s:%s %s", *artifact.Source, ss.name, artifact.Destination, localutils.TransferRate(copied, startTime))
	}
}

func (ss *SSH) Close() {
	ss.conn.Close()
}
