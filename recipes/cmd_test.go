package recipes

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

func TestRecipeCmd(t *testing.T) {
	log := logger.NewCLILogger(5, 8)
	c := Cmd{log: log}

	ls := "ls cmd*"
	c.Run([]config.Command{
		{Exec: &ls},
	})
	utils.Test().Equals(t, 3, len(c.output))
}
