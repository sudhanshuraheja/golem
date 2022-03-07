package artifacts

import (
	"github.com/sudhanshuraheja/golem/config"
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

func NewArtifact(artf config.Artifact) Artifact {
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
	return art
}
