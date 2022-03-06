package recipes

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
	"github.com/sudhanshuraheja/golem/template"
)

func TestRecipe(t *testing.T) {
	log := logger.NewCLILogger(5, 8)

	publicIP := "127.0.0.1"
	servers := []config.Server{
		{
			Name:     "one",
			HostName: []string{"one"},
			Port:     22,
			Tags:     []string{"one"},
		},
		{
			Name:     "two",
			PublicIP: &publicIP,
			HostName: []string{"two"},
			Port:     22,
			Tags:     []string{"one", "two"},
		},
		{
			Name:     "three",
			HostName: []string{"three"},
			Port:     22,
			Tags:     []string{"one", "two", "three"},
		},
	}

	match := config.Match{}
	match.Attribute = "tags"
	match.Operator = "contains"
	match.Value = "one"

	isTrue := true
	install := []string{"a", "b"}
	installNU := []string{"c", "d"}

	command1 := "ls -la @golem.key"
	commands := []string{command1}
	custom := []config.Command{
		{Exec: &command1},
		{Apt: []config.Apt{{Update: &isTrue, Install: &install, InstallNoUpgrade: &installNU}}},
	}

	licencePath := "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/LICENSE"

	recipe := config.Recipe{}
	recipe.Name = "test"
	recipe.Type = "remote"
	recipe.Match = &match
	recipe.Commands = &commands
	recipe.CustomCommands = custom
	recipe.Artifacts = []config.Artifact{
		{
			Source:      &licencePath,
			Destination: "",
		},
	}

	r := Recipe{}
	r.base = &recipe
	r.log = log

	tpl := template.Template{}
	tpl.Servers = servers
	tpl.Vars = map[string]string{
		"key": "value",
	}

	r.FindServers(servers, nil)
	utils.Test().Equals(t, 3, len(r.servers))

	match.Value = "two"
	r.FindServers(servers, nil)
	utils.Test().Equals(t, 2, len(r.servers))

	r.PrepareCommands(&tpl)
	utils.Test().Equals(t, 5, len(r.preparedCommands))
	utils.Test().Contains(t, *r.preparedCommands[1].Exec, "sudo apt-get update")
	utils.Test().Contains(t, *r.preparedCommands[2].Exec, "sudo apt-get install")

	r.PrepareArtifacts(&tpl, true)
	r.DownloadArtifacts()
	if localutils.DetectCI() {
		utils.Test().Contains(t, *r.preparedArtifacts[0].Source, "/tmp")
	} else {
		utils.Test().Contains(t, *r.preparedArtifacts[0].Source, "/var/folders/")
	}

}
