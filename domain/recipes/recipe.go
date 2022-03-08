package recipes

import (
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/natives"
	"github.com/sudhanshuraheja/golem/domain/servers"
)

type Recipe struct {
	Name      natives.String      `hcl:"name,label"`
	Type      natives.String      `hcl:"type,label"`
	Match     *servers.Match      `hcl:"match,block"`
	KeyValues kv.KeyValues        `hcl:"kv,block"`
	Artifacts artifacts.Artifacts `hcl:"artifact,block"`
	Commands  commands.Commands   `hcl:"commands"`
	Scripts   commands.Script     `hcl:"script,block"`
}
