package recipes

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
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

	command1 := "ls -la {{ .Vars.key}}"
	commands := []string{command1}
	custom := []config.Command{{Exec: &command1}}

	recipe := config.Recipe{}
	recipe.Name = "test"
	recipe.Type = "remote"
	recipe.Match = &match
	recipe.Commands = &commands
	recipe.CustomCommands = custom
	recipe.Artifacts = []config.Artifact{
		{
			Source:      "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/LICENSE",
			Destination: "",
		},
	}

	r := Recipe{}
	r.base = &recipe
	r.log = log

	template := Template{}
	template.Servers = servers
	template.Vars = map[string]string{
		"key": "value",
	}

	r.FindServers(servers)
	utils.Test().Equals(t, 3, len(r.servers))

	match.Value = "two"
	r.FindServers(servers)
	utils.Test().Equals(t, 2, len(r.servers))

	r.PrepareCommands(&template)
	utils.Test().Equals(t, 2, len(r.preparedCommands))

	r.DownloadArtifacts()
	utils.Test().Contains(t, r.base.Artifacts[0].Source, "/var/folders/")
}
