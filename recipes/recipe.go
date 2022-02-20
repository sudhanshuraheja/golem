package recipes

import (
	"fmt"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
)

func List(c *config.Config) {
	tb := log.NewTable("Name", "Match", "Artifacts", "Commands")
	for _, r := range c.Recipe {
		tb.Row(
			r.Name,
			fmt.Sprintf("%s %s %s", r.Match.Attribute, r.Match.Operator, r.Match.Value),
			len(r.Artifacts),
			len(r.Commands),
		)
	}
	// Add system defined
	tb.Row("servers", "local only", 0, 0)
	tb.Row("tflist", "local only", 0, 0)
	tb.Row("tflistall", "local only", 0, 0)
	tb.Display()
}
