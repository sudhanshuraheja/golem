package recipes

import (
	"strings"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func Servers(c *config.Config) {
	t := log.NewTable("Name", "Public IP", "Private IP", "User", "Port", "Tags", "Hostname")
	for _, s := range c.Servers {
		hostnames := utils.StringPtrValue(s.HostName, "")
		if len(hostnames) > 60 {
			hostnames = hostnames[:60]
		}
		t.Row(
			s.Name,
			utils.StringPtrValue(s.PublicIP, ""),
			utils.StringPtrValue(s.PrivateIP, ""),
			s.User,
			s.Port,
			strings.Join(s.Tags, ", "),
			hostnames,
		)
	}
	t.Display()
}
