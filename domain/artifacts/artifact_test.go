package artifacts

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/domain/vars"
)

func getArtifact(data, path, source, destination string) Artifact {
	return Artifact{
		Template: &ArtifactTemplate{
			Data: &data,
			Path: &path,
		},
		Source:      &source,
		Destination: destination,
	}
}

func TestArtifact(t *testing.T) {
	log := logger.NewCLILogger(6, 8)

	srvs := servers.Servers{}
	vrs := vars.NewVars()
	vrs.Add("key", "value")
	tpl := template.NewTemplate(srvs, *vrs, nil)

	art := getArtifact("data", "@golem.key", "source", "destination")
	err := art.TemplatePathPopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", *art.Template.Path)

	path := "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/testdata/template.tpl"
	art = getArtifact("data", path, "source", "destination")
	err = art.TemplatePathDownload(log)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, *art.Template.Path != path)

	path = "../../testdata/template.tpl"
	art = getArtifact("data", path, "source", "destination")
	err = art.TemplatePathToData()
	utils.Test().Nil(t, err)
	utils.Test().Contains(t, *art.Template.Data, "APP")

	err = art.TemplateDataPopulate(tpl)
	utils.Test().Contains(t, err.Error(), "found templates with no matches")

	vrs.Add("APP", "golem")
	err = art.TemplateDataPopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Contains(t, *art.Template.Data, "golem")

	err = art.TemplateDataToSource()
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, *art.Source != "source")

	art = getArtifact("data", "", "@golem.key", "destination")
	err = art.SourcePopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", *art.Source)

	art = getArtifact("data", "", "", "@golem.key")
	err = art.DestinationPopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", art.Destination)
}
