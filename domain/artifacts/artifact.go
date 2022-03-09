package artifacts

import (
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

func (a *Artifact) TemplatePathPopulate(tpl *template.Template) error {
	if a.Template != nil && a.Template.Path != nil {
		templatePath, err := tpl.Execute(*a.Template.Path)
		if err != nil {
			return err
		}
		a.Template.Path = &templatePath
	}
	return nil
}

func (a *Artifact) TemplatePathDownload(log *logger.CLILogger) error {
	if a.Template != nil && a.Template.Path != nil {
		if strings.HasPrefix(*a.Template.Path, "http://") || strings.HasPrefix(*a.Template.Path, "https://") {
			// Url based template
			templatePath, err := localutils.Download(log, "", *a.Template.Path)
			if err != nil {
				return err
			}
			a.Template.Path = &templatePath
		} // else File base template
	}
	return nil
}

func (a *Artifact) TemplatePathToData() error {
	if a.Template != nil && a.Template.Path != nil {
		bytes, err := os.ReadFile(*a.Template.Path)
		if err != nil {
			return err
		}
		bytesString := string(bytes)
		a.Template.Data = &bytesString
	}
	return nil
}

func (a *Artifact) TemplateDataPopulate(tpl *template.Template) error {
	if a.Template != nil && a.Template.Data != nil {
		templateData, err := tpl.Execute(*a.Template.Data)
		if err != nil {
			return err
		}
		a.Template.Data = &templateData
	}
	return nil
}

func (a *Artifact) TemplateDataToSource() error {
	if a.Template != nil && a.Template.Data != nil {
		fileName, err := localutils.FileCopy(*a.Template.Data)
		if err != nil {
			return err
		}
		a.Source = &fileName
	}
	return nil
}

func (a *Artifact) SourcePopulate(tpl *template.Template) error {
	if a.Source != nil {
		source, err := tpl.Execute(*a.Source)
		if err != nil {
			return err
		}
		a.Source = &source
	}
	return nil
}

func (a *Artifact) SourceDownload(log *logger.CLILogger) error {
	if a.Source != nil {
		if strings.HasPrefix(*a.Source, "http://") || strings.HasPrefix(*a.Source, "https://") {
			filePath, err := localutils.Download(log, "", *a.Source)
			if err != nil {
				return err
			}
			a.Source = &filePath
		}
	}
	return nil
}

func (a *Artifact) DestinationPopulate(tpl *template.Template) error {
	parsedDestination, err := tpl.Execute(a.Destination)
	if err != nil {
		return err
	}
	a.Destination = parsedDestination
	return err
}

func (a *Artifact) GetSource() string {
	source := ""
	switch {
	case a.Template != nil && a.Template.Data != nil:
		source = *a.Template.Data
	case a.Template != nil && a.Template.Path != nil:
		source = *a.Template.Path
	case a.Source != nil:
		source = *a.Source
	}
	return source
}
