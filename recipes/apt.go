package recipes

import (
	"github.com/sudhanshuraheja/golem/config"
)

func AptUpdate(c *config.Config) {
	for _, s := range c.Servers.Server {
		SSHRun(&s, []string{
			"apt-get update",
		})
	}
}
