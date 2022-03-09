package recipes

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/domain/vars"
)

func TestRecipes(t *testing.T) {
	log := logger.NewCLILogger(6, 8)

	source := "source"
	rcps := Recipes{}
	utils.Test().Equals(t, 0, len(rcps))
	rcps.Append(Recipe{
		Name:     "test1",
		Type:     "remote",
		Commands: &[]commands.Command{"ls -la"},
		Artifacts: []*artifacts.Artifact{
			{Source: &source, Destination: "destination"},
		},
	})
	utils.Test().Equals(t, 1, len(rcps))
	rcps2 := Recipes{}
	rcps.Merge(rcps2)
	utils.Test().Equals(t, 1, len(rcps))

	vrs := vars.NewVars()
	vrs.Add("key", "value")
	tpl := template.NewTemplate([]servers.Server{}, *vrs, nil)

	rcps.Display(log, tpl, "test1")

	rcpSelected, err := rcps.Search("test1")
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "test1", rcpSelected.Name)
}
