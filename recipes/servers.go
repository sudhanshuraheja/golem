package recipes

import (
	"strings"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func Servers(c *config.Config) {
	t := log.NewTable("Name", "Public IP", "Private IP", "Hostname", "User", "Port", "Tags")
	for _, s := range c.Servers {
		t.Row(
			s.Name,
			utils.StringPtrValue(s.PublicIP, ""),
			utils.StringPtrValue(s.PrivateIP, ""),
			utils.StringPtrValue(s.HostName, ""),
			s.User,
			s.Port,
			strings.Join(s.Tags, ","),
		)
	}
	t.Display()
}
