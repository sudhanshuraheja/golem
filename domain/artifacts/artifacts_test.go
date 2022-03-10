package artifacts

import (
	"os"
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/domain/vars"
)

func TestArtifacts(t *testing.T) {
	art := NewArtifact("@golem.APP", "", "", "destination")
	art2 := NewArtifact("", "../../testdata/template.tpl", "", "destination2")
	art3 := NewArtifact("", "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/testdata/template.tpl", "", "destination2")
	art4 := NewArtifact("", "", "../../testdata/template.tpl", "destination2")
	art5 := NewArtifact("", "", "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/testdata/template.tpl", "destination2")

	arts := Artifacts{}
	arts.Append(*art)
	arts.Append(*art2)
	arts.Append(*art3)
	utils.Test().Equals(t, 3, len(arts))

	arts2 := Artifacts{}
	arts2.Append(*art4)
	arts2.Append(*art5)
	utils.Test().Equals(t, 2, len(arts2))

	arts.Merge(arts2)
	utils.Test().Equals(t, 5, len(arts))

	log := logger.NewCLILogger(6, 8)
	srvs := servers.Servers{}
	vrs := vars.NewVars()
	vrs.Add("APP", "golem")
	tpl := template.NewTemplate(srvs, *vrs, nil)

	arts.PrepareForExecution(log, tpl)

	output := map[int]string{
		0: "golem",
		1: "golem-golem",
		2: "golem-golem",
		3: "{{ .Vars.APP }}-@golem.APP",
		4: "{{ .Vars.APP }}-@golem.APP",
	}

	for i, a := range arts {
		bytes, err := os.ReadFile(*a.Source)
		utils.Test().Nil(t, err)
		utils.Test().Equals(t, output[i], string(bytes))
	}
}
