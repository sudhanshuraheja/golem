package recipes

import (
	"fmt"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/template"
)

type Recipes []Recipe

func (r *Recipes) Append(rcp Recipe) {
	if r != nil {
		*r = append(*r, rcp)
	}
}

func (r *Recipes) Merge(rcps Recipes) {
	if r != nil {
		*r = append(*r, rcps...)
	}
}

func (r *Recipes) Display(log *logger.CLILogger, tpl *template.Template, query string) {
	log.Announce("").Msgf("list of all available recipes")

	// Add system defined
	log.Info("system").Msgf("%s", logger.Cyan("recipes"))
	log.Info("system").Msgf("%s", logger.Cyan("servers"))

	for _, recipe := range *r {
		recipe.Display(log, tpl, query)
	}
}

func (r *Recipes) Search(name string) (*Recipe, error) {
	for _, rcp := range *r {
		if string(rcp.Name) == name {
			return &rcp, nil
		}
	}
	return nil, fmt.Errorf("could not find this recipe")
}
