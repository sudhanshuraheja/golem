package recipes

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/keyvalue"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/domain/vars"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

func TestRecipe(t *testing.T) {
	log := logger.NewCLILogger(6, 8)
	vars := make(map[string]string)
	store := kv.NewStore(log)
	tpl := template.NewTemplate(servers.Servers{}, vars, store)

	trueValue := true
	cmd := commands.NewCommand("ls -la script2")
	kval := keyvalue.KeyValue{Path: "test.key", Value: "value"}
	keyvalues := []*keyvalue.KeyValue{&kval}

	script := commands.Script{
		Apt:      []commands.Apt{{Update: &trueValue}},
		Commands: &[]commands.Command{"ls -la script"},
		Command:  &cmd,
	}

	r := Recipe{
		Name:      "test1",
		Type:      "local",
		Match:     servers.NewMatch("name", "=", "test"),
		KeyValues: keyvalues,
		Scripts:   []*commands.Script{&script},
		Artifacts: []*artifacts.Artifact{
			{Source: localutils.StrPtr("source"), Destination: "destination"},
		},
		Commands: &[]commands.Command{"ls -la"},
	}

	r.Display(log, tpl, "")
	r.AskPermission(log)

	err := r.Prepare(log, store)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 1, len(r.KeyValues))
	utils.Test().Equals(t, 1, len(r.Scripts))
	utils.Test().Equals(t, 2, len(r.Artifacts))
	utils.Test().Equals(t, 6, len(*r.Commands))

	err = r.PrepareForExecution(log, tpl, store)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 1, len(r.KeyValues))
	utils.Test().Equals(t, 1, len(r.Scripts))
	utils.Test().Equals(t, 2, len(r.Artifacts))
	utils.Test().Equals(t, 6, len(*r.Commands))

	err = store.Delete("test.key")
	utils.Test().Nil(t, err)
}

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
