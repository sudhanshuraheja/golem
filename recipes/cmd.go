package recipes

import (
	"os"
	"path/filepath"
	"time"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/cmd"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Cmd struct {
	log    *logger.CLILogger
	output []cmd.Out
}

func (c *Cmd) Run(commands []string) {
	name := "local"
	c.output = []cmd.Out{}
	cm := cmd.NewCmd(name)

	wait := make(chan bool)
	go func(log *logger.CLILogger, wait chan bool) {
		for {
			select {
			case stdout := <-cm.Stdout:
				c.output = append(c.output, stdout)
				if stdout.Message != "" {
					c.log.Debug(stdout.Name).Msgf("%s", stdout.Message)
				}
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-cm.Stderr:
				c.output = append(c.output, stderr)
				if stderr.Message != "" {
					c.log.Error(stderr.Name).Msgf("%s", stderr.Message)
				}
				if stderr.Completed {
					wait <- true
				}
			}
		}
	}(c.log, wait)

	for _, command := range commands {
		startTime := time.Now()
		c.log.Highlight(name).Msgf("$ %s", command)
		err := cm.Run(command)
		if err != nil {
			c.log.Error(name).Msgf("error in running command <%s>: %v", command, err)
			continue
		}
		<-wait
		c.log.Success(name).Msgf("$ %s %s", command, localutils.TimeInSecs(startTime))
	}
}

func (c *Cmd) Upload(artifacts []config.Artifact) {
	name := "local"
	for _, artifact := range artifacts {
		startTime := time.Now()
		err := os.MkdirAll(filepath.Dir(artifact.Destination), os.ModePerm)
		if err != nil {
			c.log.Error(name).Msgf("could not create directory <%s>: %v", filepath.Dir(artifact.Destination), err)
			continue
		}
		err = os.Rename(*artifact.Source, artifact.Destination)
		if err != nil {
			c.log.Error(name).Msgf("error in moving from <%s> to <%s>: %v", *artifact.Source, artifact.Destination, err)
			continue
		}
		c.log.Success(name).Msgf("%s %s %s %s %s", logger.Cyan("Moved"), *artifact.Source, logger.Cyan("to"), artifact.Destination, localutils.TimeInSecs(startTime))
	}
}
