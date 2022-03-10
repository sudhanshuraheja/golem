package artifacts

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/domain/vars"
)

func TestArtifact(t *testing.T) {
	log := logger.NewCLILogger(6, 8)

	srvs := servers.Servers{}
	vrs := vars.NewVars()
	vrs.Add("key", "value")
	tpl := template.NewTemplate(srvs, *vrs, nil)

	art := NewArtifact("data", "@golem.key", "source", "destination")
	err := art.TemplatePathPopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", *art.Template.Path)
	artSource := art.GetSource()
	utils.Test().Equals(t, "data", artSource)

	path := "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/testdata/template.tpl"
	art = NewArtifact("data", path, "source", "destination")
	err = art.TemplatePathDownload(log)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, *art.Template.Path != path)

	path = "../../testdata/template.tpl"
	art = NewArtifact("data", path, "source", "destination")
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

	art = NewArtifact("data", "", "@golem.key", "destination")
	err = art.SourcePopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", *art.Source)

	art = NewArtifact("data", "", "", "@golem.key")
	err = art.DestinationPopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", art.Destination)

	// Without template
	art = NewArtifact("", "", "source", "destination")
	err = art.TemplatePathPopulate(tpl)
	utils.Test().Nil(t, err)
	artSource = art.GetSource()
	utils.Test().Equals(t, "source", artSource)

	art = NewArtifact("", "", "source", "destination")
	err = art.TemplatePathDownload(log)
	utils.Test().Nil(t, err)

	art = NewArtifact("", "", "source", "destination")
	err = art.TemplatePathToData()
	utils.Test().Nil(t, err)

	err = art.TemplateDataPopulate(tpl)
	utils.Test().Nil(t, err)

	vrs.Add("APP", "golem")
	err = art.TemplateDataPopulate(tpl)
	utils.Test().Nil(t, err)

	err = art.TemplateDataToSource()
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "source", *art.Source)

	art = NewArtifact("", "", "@golem.key", "destination")
	err = art.SourcePopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", *art.Source)

	art = NewArtifact("", "", "", "@golem.key")
	err = art.DestinationPopulate(tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "value", art.Destination)
}
