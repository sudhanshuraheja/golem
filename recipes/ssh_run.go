package recipes

import (
	"time"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/ssh"
)

func SSHRun(c *config.Config, commands []string) {
	for _, s := range c.Servers.Server {
		startTime := time.Now()
		log.Announcef("%s | running <apt-update>", s.Name)

		var host string
		switch {
		case s.PublicIP != nil:
			host = *s.PublicIP
		case s.HostName != nil:
			host = *s.HostName
		default:
			log.Errorf("%s | could not find a public ip or a hostname in config", s.Name)
		}

		conn, err := ssh.NewSSHConnection(s.Name, s.User, host, s.Port, "")
		if err != nil {
			log.Errorf("%s | could not ssh to host: %v", s.Name, err)
		}

		wait := make(chan bool)

		go func(wait chan bool) {
			for {
				select {
				case stdout := <-conn.Stdout:
					log.Infof("%s | %s", stdout.Name, stdout.Message)
					if stdout.Completed {
						wait <- true
					}
				case stderr := <-conn.Stderr:
					log.Infof("%s | %s", stderr.Name, stderr.Message)
				}
			}
		}(wait)

		for _, cmd := range commands {
			status, err := conn.Run(cmd)
			if err != nil {
				log.Errorf("%s | error in running command <%s>: %v", s.Name, cmd, err)
			}
			<-wait
			if status > 0 {
				log.Errorf("%s | command <%s> failed with status: %d", s.Name, cmd, status)
			} else {
				log.Successf("%s | command <%s> ended successfully in %s", s.Name, cmd, time.Since(startTime))
			}
		}

		conn.Close()
	}
}
