package recipes

import (
	"fmt"
	"time"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/ssh"
)

type SSH struct {
	conn *ssh.Connection
	name string
}

type SSHJob struct {
	Recipe config.Recipe
	Server config.Server
}

func (ss *SSH) Connect(s *config.Server) error {
	var host string
	switch {
	case s.PublicIP != nil:
		host = *s.PublicIP
	case s.HostName != nil:
		host = *s.HostName
	default:
		return fmt.Errorf("could not find a public ip or a hostname in config")
	}

	startTime := time.Now()
	conn, err := ssh.NewSSHConnection(s.Name, s.User, host, s.Port, "")
	if err != nil {
		return fmt.Errorf("could not ssh to host: %v", err)
	}
	log.MinorSuccessf("%s | connected via SSH in %s", s.Name, time.Since(startTime))
	ss.conn = conn
	ss.name = s.Name
	return nil
}

func (ss *SSH) Run(commands []string) {
	wait := make(chan bool)
	go func(wait chan bool) {
		for {
			select {
			case stdout := <-ss.conn.Stdout:
				log.Infof("%s | %s", stdout.Name, stdout.Message)
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-ss.conn.Stderr:
				log.Infof("%s | %s", stderr.Name, stderr.Message)
			}
		}
	}(wait)

	for _, cmd := range commands {
		log.Announcef("%s | running <%s>", ss.name, cmd)
		startTime := time.Now()
		status, err := ss.conn.Run(cmd)
		if err != nil {
			log.Errorf("%s | error in running command <%s>: %v", ss.name, cmd, err)
			continue
		}
		<-wait
		if status > 0 {
			log.Errorf("%s | command <%s> failed with status: %d", ss.name, cmd, status)
			continue
		}
		log.Successf("%s | command <%s> ended successfully in %s", ss.name, cmd, time.Since(startTime))
	}
}

func (ss *SSH) Upload(artifacts []config.Artifact) {
	for _, artifact := range artifacts {
		startTime := time.Now()
		copied, err := ss.conn.Upload(artifact.Source, artifact.Destination)
		if err != nil {
			log.Errorf("%s | failed to upload local:<%s> to remote:<%s>: %v", ss.name, artifact.Source, artifact.Destination, err)
			continue
		}
		log.Successf("%s | successfully uploaded %d bytes from local:<%s> to remote:<%s> in %s", ss.name, copied, artifact.Source, artifact.Destination, time.Since(startTime))
	}
}

func (ss *SSH) Close() {
	ss.conn.Close()
}
