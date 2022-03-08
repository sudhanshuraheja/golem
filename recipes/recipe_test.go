package recipes

import (
	"strings"
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/commands"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/plugins"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
)

func TestRecipe(t *testing.T) {
	log := logger.NewCLILogger(5, 8)

	conf := config.NewConfig("../testdata/sample.hcl")
	srvs := conf.Servers
	tpl := template.NewTemplate(srvs, *conf.Vars, nil)

	recipes := []Recipe{}
	for _, crcp := range conf.Recipes {
		rcp := NewRecipe(log, tpl)
		rcp.Name = crcp.Name
		rcp.OfType = crcp.Type
		if crcp.Match != nil {
			rcp.Match = servers.NewMatch(crcp.Match.Attribute, crcp.Match.Operator, crcp.Match.Value)
		}
		for _, keyValue := range crcp.KeyValues {
			rcp.KV[keyValue.Path] = keyValue.Value
		}
		for _, cmd := range crcp.CustomCommands {
			if cmd.Exec != nil {
				parsedCmd, err := tpl.Execute(*cmd.Exec)
				utils.Test().Nil(t, err)
				parsedCmd = strings.TrimSuffix(parsedCmd, "\n")
				rcp.AddCommand(commands.Command{Exec: parsedCmd})
			}

			apt := plugins.NewAPT()
			cmds, artfs := apt.Prepare(cmd.Apt)
			for _, cmd := range cmds {
				rcp.AddCommand(cmd)
			}
			rcp.PrepareArtifacts(artfs, true)
		}
		if crcp.Commands != nil {
			for _, cmd := range *crcp.Commands {
				parsedCmd, err := tpl.Execute(cmd)
				utils.Test().Nil(t, err)
				parsedCmd = strings.TrimSuffix(parsedCmd, "\n")
				rcp.AddCommand(commands.Command{Exec: parsedCmd})
			}
		}
		artfs := []artifacts.Artifact{}
		for _, art := range crcp.Artifacts {
			artfs = append(artfs, art)
		}
		rcp.PrepareArtifacts(artfs, true)
		recipes = append(recipes, *rcp)
	}

	recipes[0].findServers(srvs)
	utils.Test().Equals(t, 1, len(recipes[0].servers))
}
