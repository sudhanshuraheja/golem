package recipes

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
)

func TestRecipeCmd(t *testing.T) {
	log := logger.NewCLILogger(5, 8)
	c := Cmd{log: log}
	c.Run([]string{"ls"})
	utils.Test().Equals(t, 12, len(c.output))
}
