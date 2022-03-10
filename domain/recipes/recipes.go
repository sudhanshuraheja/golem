package recipes

import (
	"fmt"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/kv"
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
		if rcp.Name == name {
			return &rcp, nil
		}
	}
	return nil, fmt.Errorf("could not find this recipe")
}

func (r *Recipes) Prepare(log *logger.CLILogger, store *kv.Store) error {
	if r != nil {
		for i, rcp := range *r {
			err := rcp.Prepare(log, store)
			if err != nil {
				return err
			}
			(*r)[i] = rcp
		}
	}
	return nil
}

func (r *Recipes) PrepareForExecution(log *logger.CLILogger, tpl *template.Template) error {
	if r != nil {
		for i, rcp := range *r {
			err := rcp.PrepareForExecution(log, tpl)
			if err != nil {
				return err
			}
			(*r)[i] = rcp
		}
	}
	return nil
}
