package kitchen

import (
	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/artifacts"
	"github.com/sudhanshuraheja/golem/commands"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/kv"
	"github.com/sudhanshuraheja/golem/plugins"
	"github.com/sudhanshuraheja/golem/recipes"
	"github.com/sudhanshuraheja/golem/servers"
	"github.com/sudhanshuraheja/golem/template"
)

type Recipes struct {
	conf    *config.Config
	log     *logger.CLILogger
	tpl     *template.Template
	kv      *kv.KV
	Servers []servers.Server
	Recipes []recipes.Recipe
}

func NewRecipes(conf *config.Config, log *logger.CLILogger) *Recipes {
	r := Recipes{
		conf: conf,
		log:  log,
	}
	r.Servers = servers.NewServers(conf.Servers)
	r.kv = kv.NewKV(log)

	vars := map[string]string{}
	if conf.Vars != nil {
		vars = *conf.Vars
	}
	r.tpl = template.NewTemplate(r.Servers, vars, r.kv)

	for _, crcp := range conf.Recipes {
		rcp := recipes.NewRecipe(log, r.tpl)
		rcp.Name = crcp.Name
		rcp.OfType = crcp.Type
		if crcp.Match != nil {
			rcp.Match = servers.NewMatch(crcp.Match.Attribute, crcp.Match.Operator, crcp.Match.Value)

			var err error
			rcp.Match.Value, err = r.tpl.Execute(rcp.Match.Value)
			if err != nil {
				r.log.Error(rcp.Name).Msgf("could not parse template %s: %v", rcp.Match.Value, err)
			}
		}
		for _, keyValue := range crcp.KeyValues {
			rcp.KV[keyValue.Path] = keyValue.Value
		}
		for _, cmd := range crcp.CustomCommands {
			if cmd.Exec != nil {
				rcp.AddCommand(commands.Command{Exec: *cmd.Exec})
			}

			apt := plugins.NewAPT()
			cmds, artfs := apt.Prepare(cmd.Apt)
			for _, cmd := range cmds {
				rcp.AddCommand(cmd)
			}
			for _, artf := range artfs {
				rcp.AddArtifact(artf)
			}
		}
		if crcp.Commands != nil {
			for _, cmd := range *crcp.Commands {
				rcp.AddCommand(commands.Command{Exec: cmd})
			}
		}
		for _, art := range crcp.Artifacts {
			rcp.AddArtifact(artifacts.NewArtifact(art))
		}

		r.Recipes = append(r.Recipes, *rcp)
	}

	return &r
}

func (r *Recipes) ListRecipes(query string) {
	r.log.Announce("").Msgf("list of all available recipes")

	// Add system defined
	r.log.Info("system").Msgf("%s", logger.Cyan("recipes"))
	r.log.Info("system").Msgf("%s", logger.Cyan("servers"))

	for _, recipe := range r.Recipes {
		recipe.Display(query)
	}
}

func (r *Recipes) ListServers(query string) {
	r.log.Announce("").Msgf("list of all connected servers")
	for _, s := range r.Servers {
		s.Display(r.log, query)
	}
}

func (r *Recipes) KV(path, action string) {
	if path == "" || path == "list" {
		r.kv.Display(r.log, action)
		return
	}

	switch action {
	case "set":
		r.kv.SetUserValue(path)
	case "rand32":
		r.kv.SetValue(path, action)
	case "delete":
		r.kv.DeleteValue(path)
	default:
		r.kv.GetValue(path)
	}

	_ = r.kv.Close()
}

func (r *Recipes) Run(name string) {
	var recipe *recipes.Recipe
	for i, rcp := range r.Recipes {
		if rcp.Name == name {
			recipe = &r.Recipes[i]
		}
	}

	if recipe.Name == "" {
		r.log.Error(name).Msgf("the recipe %s was not found in '~/.golem/' or '.'", logger.Cyan(name))
		return
	}

	maxParallelProcesses := 4
	if r.conf.MaxParallelProcesses != nil {
		maxParallelProcesses = *r.conf.MaxParallelProcesses
	}

	recipe.Execute(r.Servers, r.kv, maxParallelProcesses)
}
