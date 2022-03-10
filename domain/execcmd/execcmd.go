package execcmd

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/pkg/cmd"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type ExecCmd struct {
	log    *logger.CLILogger
	output []cmd.Out
	mu     sync.RWMutex
}

func NewExecCmd(log *logger.CLILogger) *ExecCmd {
	return &ExecCmd{
		log: log,
	}
}

func (c *ExecCmd) Start(cmds *[]commands.Command, artfs []*artifacts.Artifact) {
	if artfs != nil {
		c.Upload(artfs)
	}
	if cmds != nil {
		c.Run(*cmds)
	}
}

func (c *ExecCmd) OutputAppend(out cmd.Out) {
	c.mu.Lock()
	c.output = append(c.output, out)
	c.mu.Unlock()
}

func (c *ExecCmd) Run(cmds commands.Commands) {
	c.output = []cmd.Out{}

	name := "local"
	shell := cmd.NewCmd(name)

	wait := make(chan bool)
	go func(wait chan bool) {
		for {
			select {
			case stdout := <-shell.Stdout:
				c.OutputAppend(stdout)
				if stdout.Message != "" {
					c.log.Debug(stdout.Name).Msgf("%s", stdout.Message)
				}
				if stdout.Completed {
					wait <- true
				}
			case stderr := <-shell.Stderr:
				c.OutputAppend(stderr)
				if stderr.Message != "" {
					c.log.Error(stderr.Name).Msgf("%s", stderr.Message)
				}
				if stderr.Completed {
					wait <- true
				}
			}
		}
	}(wait)

	for _, cmd := range cmds {
		if cmd == "" {
			continue
		}
		startTime := time.Now()
		c.log.Highlight(name).Msgf("$ %s", cmd)
		err := shell.Run(string(cmd))
		if err != nil {
			c.log.Error(name).Msgf("error in running command <%s>: %v", cmd, err)
			continue
		}
		<-wait
		c.log.Success(name).Msgf("$ %s %s", cmd, localutils.TimeInSecs(startTime))
	}
}

func (c *ExecCmd) Upload(artfs []*artifacts.Artifact) {
	name := "local"
	for _, artf := range artfs {
		startTime := time.Now()
		err := os.MkdirAll(filepath.Dir(artf.Destination), os.ModePerm)
		if err != nil {
			c.log.Error(name).Msgf("could not create directory <%s>: %v", filepath.Dir(artf.Destination), err)
			continue
		}
		if artf.Source == nil {
			c.log.Error(name).Msgf("source has not been setup yet for destination %s", artf.Destination)
			continue
		}
		err = os.Rename(*artf.Source, artf.Destination)
		if err != nil {
			c.log.Error(name).Msgf("error in moving from <%s> to <%s>: %v", artf.Source, artf.Destination, err)
			continue
		}
		c.log.Success(name).Msgf(
			"%s %s %s %s %s",
			logger.Cyan("Moved"),
			localutils.TinyString(*artf.Source, 50),
			logger.Cyan("to"),
			artf.Destination,
			localutils.TimeInSecs(startTime))
	}
}
