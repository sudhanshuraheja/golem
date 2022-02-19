package recipes

import (
	"time"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/ssh"
)

func SSHRun(c *config.Config, commands []string) {
	for _, s := range c.Servers.Server {
		startTime := time.Now()
		Info().Println(s.Name, "| Running apt-update")

		var host string
		switch {
		case s.PublicIP != nil:
			host = *s.PublicIP
		case s.HostName != nil:
			host = *s.HostName
		default:
			Errors().Println(s.Name, "| could not find a public ip or a hostname")
		}

		conn, err := ssh.NewSSHConnection(s.Name, s.User, host, s.Port, "")
		if err != nil {
			Errors().Println(s.Name, "| Could not ssh to host", err)
		}

		wait := make(chan bool)

		go func(wait chan bool) {
			for {
				select {
				case stdout := <-conn.Stdout:
					Progress().Println(stdout.Name, "|", stdout.Message)
					if stdout.Completed {
						wait <- true
					}
				case stderr := <-conn.Stderr:
					Progress().Println(stderr.Name, "|", stderr.Message)
				}
			}
		}(wait)

		for _, cmd := range commands {
			status, err := conn.Run(cmd)
			if err != nil {
				Errors().Println(s.Name, "| Error during running", cmd, err)
			}
			<-wait
			if status > 0 {
				Errors().Println(s.Name, "| Command", cmd, "ended with status", status)
			} else {
				Success().Println(s.Name, "| Command ended in", time.Since(startTime))
			}
		}

		conn.Close()
	}
}
