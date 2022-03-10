package execssh

import (
	"fmt"
	"time"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
	"github.com/sudhanshuraheja/golem/pkg/ssh"
)

type SSH struct {
	conn   ssh.Connection
	log    *logger.CLILogger
	name   string
	output []ssh.Out
}

func (ss *SSH) Connect(s *servers.Server) error {
	host, err := s.GetHostName()
	if err != nil {
		return err
	}

	ss.name = s.Name
	startTime := time.Now()
	conn, err := ssh.NewSSHConnection(s.Name, s.User, host, s.Port, "")
	if err != nil {
		return fmt.Errorf("could not ssh to host: %v", err)
	}

	ss.log.Info(s.Name).Msgf("connected via SSH %s", localutils.TimeInSecs(startTime))
	ss.conn = conn

	return nil
}

func (ss *SSH) Run(cmds commands.Commands) {
	ss.output = []ssh.Out{}
	wait := make(chan bool)

	go func(wait chan bool) {
		stdoutCh := ss.conn.Stdout()
		stderrCh := ss.conn.Stderr()
		for {
			select {
			case stdout := <-stdoutCh:
				ss.output = append(ss.output, stdout)
				if stdout.Message != "" {
					ss.log.Debug(stdout.Name).Msgf("%s", stdout.Message)
				}
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-stderrCh:
				ss.output = append(ss.output, stderr)
				if stderr.Message != "" {
					ss.log.Error(stderr.Name).Msgf("%s", stderr.Message)
				}
				if stderr.Completed {
					wait <- true
				}
			}
		}
	}(wait)

	for _, cmd := range cmds {
		startTime := time.Now()
		ss.log.Highlight(ss.name).Msgf("$ %s", cmd)
		status, err := ss.conn.Run(string(cmd))
		if err != nil {
			ss.log.Error(ss.name).Msgf("error in running command <%s>: %v", cmd, err)
			continue
		}
		<-wait
		if status > 0 {
			ss.log.Error(ss.name).Msgf("command <%s> failed with status: %d", cmd, status)
			continue
		}
		ss.log.Success(ss.name).Msgf("$ %s %s", cmd, localutils.TimeInSecs(startTime))
	}
}

func (ss *SSH) Upload(artfs []*artifacts.Artifact) {
	for _, artf := range artfs {
		if artf.Source == nil {
			ss.log.Error(ss.name).Msgf("source does not exist for destination %s", artf.Destination)
			return
		}

		startTime := time.Now()
		ss.log.Info(ss.name).Msgf(
			"%s %s %s %s:%s",
			logger.Cyan("uploading"),
			localutils.TinyString(*artf.Source, 50),
			logger.Cyan("to"),
			ss.name,
			localutils.TinyString(artf.Destination, 50),
		)
		copied, err := ss.conn.Upload(*artf.Source, artf.Destination)
		if err != nil {
			ss.log.Error(ss.name).Msgf(
				"failed to upload local:<%s> to remote:<%s>: %v",
				artf.Source,
				artf.Destination,
				err,
			)
			continue
		}
		ss.log.Success(ss.name).Msgf(
			"uploaded %s to %s:%s %s",
			localutils.TinyString(*artf.Source, 50),
			ss.name,
			artf.Destination,
			localutils.TransferRate(copied, startTime),
		)
	}
}

func (ss *SSH) Close() {
	ss.conn.Close()
}
