package execcmd

import (
	"testing"
	"time"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/commands"
)

func TestRecipeCmd(t *testing.T) {
	log := logger.NewCLILogger(5, 8)
	c := ExecCmd{log: log}

	c.Run(commands.Commands{
		commands.NewCommand("ls execcmd*"),
	})

	time.Sleep(time.Second)
	c.mu.Lock()
	utils.Test().Equals(t, 3, len(c.output))
	c.mu.Unlock()
}
