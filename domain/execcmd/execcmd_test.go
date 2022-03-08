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
	utils.Test().Equals(t, 3, len(c.output))
}
