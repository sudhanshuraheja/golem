package execcmd

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/commands"
)

func TestRecipeCmd(t *testing.T) {
	log := logger.NewCLILogger(5, 8)
	c := ExecCmd{log: log}

	c.Run(commands.Commands{
		commands.NewCommand("ls cmd*"),
	})

	c.mu.Lock()
	utils.Test().Equals(t, 2, len(c.output))
	c.mu.Unlock()
}
