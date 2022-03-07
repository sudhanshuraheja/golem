package artifacts

import (
	"fmt"
	"os"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
	"github.com/sudhanshuraheja/golem/template"
)

type Artifact struct {
	Template    Template
	Source      string
	Destination string
}

type Template struct {
	Data string
	Path string
}

func NewArtifact(artf config.Artifact) *Artifact {
	art := Artifact{}
	art.Template = Template{}
	if artf.Template != nil {
		if artf.Template.Data != nil {
			art.Template.Data = *artf.Template.Data
		}
		if artf.Template.Path != nil {
			art.Template.Path = *artf.Template.Path
		}
	}
	if artf.Source != nil {
		art.Source = *artf.Source
	}
	art.Destination = artf.Destination
	return &art
}

func (a *Artifact) HandlePath(log *logger.CLILogger, tpl *template.Template) error {
	var err error
	if a.Template.Path != "" {
		a.Template.Path, err = tpl.Execute(a.Template.Path)
		if err != nil {
			return err
		}

		if strings.HasPrefix(a.Template.Path, "http://") || strings.HasPrefix(a.Template.Path, "https://") {
			// Url based template
			a.Template.Path, err = localutils.Download(log, "", a.Template.Path)
			if err != nil {
				return err
			}
		} // else File base template

		bytes, err := os.ReadFile(a.Template.Path)
		if err != nil {
			return err
		}
		bytesString := string(bytes)
		a.Template.Data = bytesString
	}
	return err
}

func (a *Artifact) HandleData(tpl *template.Template, dryrun bool) error {
	var err error
	if a.Template.Data != "" {
		a.Template.Data, err = tpl.Execute(a.Template.Data)
		if err != nil {
			return err
		}

		if !dryrun {
			fileName, err := localutils.FileCopy(a.Template.Data)
			if err != nil {
				return err
			}
			a.Source = fileName
		}
	}
	return err
}

func (a *Artifact) HandleSource(tpl *template.Template) error {
	var err error
	if a.Source != "" {
		a.Source, err = tpl.Execute(a.Source)
		if err != nil {
			return err
		}
	}
	return err
}

func (a *Artifact) HandleDestination(tpl *template.Template) error {
	var err error
	a.Destination, err = tpl.Execute(a.Destination)
	return err
}

func (a *Artifact) Download(log *logger.CLILogger) (string, error) {
	if a.Source != "" {
		if strings.HasPrefix(a.Source, "http://") || strings.HasPrefix(a.Source, "https://") {
			filePath, err := localutils.Download(log, "", a.Source)
			if err != nil {
				return "", err
			}
			return filePath, nil
		}
	}
	return "", fmt.Errorf("source has not been set up yet")
}
