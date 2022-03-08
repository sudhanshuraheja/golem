package artifacts

import (
	"fmt"
	"os"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Artifact struct {
	Template    *ArtifactTemplate `hcl:"template,block"`
	Source      *string           `hcl:"source"`
	Destination string            `hcl:"destination"`
}

type ArtifactTemplate struct {
	Data *string `hcl:"data"`
	Path *string `hcl:"path"`
}

func (a *Artifact) HandlePath(log *logger.CLILogger, tpl *template.Template) error {
	if a.Template.Path != nil {
		templatePath, err := tpl.Execute(*a.Template.Path)
		if err != nil {
			return err
		}
		a.Template.Path = &templatePath

		if strings.HasPrefix(*a.Template.Path, "http://") || strings.HasPrefix(*a.Template.Path, "https://") {
			// Url based template
			templatePath, err = localutils.Download(log, "", *a.Template.Path)
			if err != nil {
				return err
			}
			a.Template.Path = &templatePath
		} // else File base template

		bytes, err := os.ReadFile(*a.Template.Path)
		if err != nil {
			return err
		}
		bytesString := string(bytes)
		a.Template.Data = &bytesString
	}
	return nil
}

func (a *Artifact) HandleData(tpl *template.Template, dryrun bool) error {
	if a.Template.Data != nil {
		templateData, err := tpl.Execute(*a.Template.Data)
		if err != nil {
			return err
		}
		a.Template.Data = &templateData

		if !dryrun {
			fileName, err := localutils.FileCopy(*a.Template.Data)
			if err != nil {
				return err
			}
			a.Source = &fileName
		}
	}
	return nil
}

func (a *Artifact) HandleSource(tpl *template.Template) error {
	if a.Source != nil {
		source, err := tpl.Execute(*a.Source)
		if err != nil {
			return err
		}
		a.Source = &source
	}
	return nil
}

func (a *Artifact) HandleDestination(tpl *template.Template) error {
	var err error
	a.Destination, err = tpl.Execute(a.Destination)
	return err
}

func (a *Artifact) Download(log *logger.CLILogger) (string, error) {
	if a.Source != nil {
		if strings.HasPrefix(*a.Source, "http://") || strings.HasPrefix(*a.Source, "https://") {
			filePath, err := localutils.Download(log, "", *a.Source)
			if err != nil {
				return "", err
			}
			return filePath, nil
		}
	}
	return "", fmt.Errorf("source has not been set up yet")
}
