package recipes

import (
	"github.com/sudhanshuraheja/golem/config"
)

func AptUpdate(c *config.Config) {
	SSHRun(c, []string{
		"apt-get update",
	})
}
