package artifacts

import (
	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/template"
)

type Artifacts []Artifact

func (a *Artifacts) Append(art Artifact) {
	*a = append(*a, art)
}

func (a *Artifacts) Merge(arts Artifacts) {
	*a = append(*a, arts...)
}

func (a *Artifacts) Prepare(log *logger.CLILogger, tpl *template.Template, dryrun bool) {
	for _, art := range *a {
		err := art.TemplatePathPopulate(tpl)
		if err != nil {
			log.Error("").Msgf("could not populate template path: %v", err)
			continue
		}

		err = art.TemplatePathDownload(log)
		if err != nil {
			log.Error("").Msgf("coult not download template path: %v", err)
			continue
		}

		err = art.TemplatePathToData()
		if err != nil {
			log.Error("").Msgf("coult not move to template data: %v", err)
			continue
		}

		err = art.TemplateDataPopulate(tpl)
		if err != nil {
			log.Error("").Msgf("coult not populate template data: %v", err)
			continue
		}

		err = art.TemplateDataToSource()
		if err != nil {
			log.Error("").Msgf("coult not move to source: %v", err)
			continue
		}

		err = art.SourcePopulate(tpl)
		if err != nil {
			log.Error("").Msgf("coult not populate source: %v", err)
			continue
		}

		err = art.SourceDownload(log)
		if err != nil {
			log.Error("").Msgf("coult not download source: %v", err)
			continue
		}

		err = art.DestinationPopulate(tpl)
		if err != nil {
			log.Error("").Msgf("coult not populate destination: %v", err)
			continue
		}
	}
}
