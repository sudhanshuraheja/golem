package artifacts

import (
	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/template"
)

type Artifacts []*Artifact

func (a *Artifacts) Append(art Artifact) {
	*a = append(*a, &art)
}

func (a *Artifacts) Merge(arts Artifacts) {
	*a = append(*a, arts...)
}

func (a *Artifacts) PrepareForExecution(log *logger.CLILogger, tpl *template.Template) {
	for i, art := range *a {
		err := art.PrepareForExecution(log, tpl)
		if err != nil {
			log.Error("").Msgf("%v", err)
		}
		(*a)[i] = art
	}
}

func (a *Artifacts) ToPointerArray() []*Artifact {
	aa := []*Artifact{}
	if a != nil {
		for _, art := range *a {
			aa = append(aa, art)
		}
	}
	return aa
}
