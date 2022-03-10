package golem

import (
	"os"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/template"
)

type Golem struct {
	cli   *Config
	conf  *config.Config
	log   *logger.CLILogger
	tpl   *template.Template
	store *kv.Store
}

func NewGolem(conf *Config) {
	g := Golem{}
	g.cli = conf
	g.log = logger.NewCLILogger(6, 12)
	g.conf = &config.Config{}

	files, err := conf.Detect(g.log)
	if err != nil {
		g.log.Fatal("golem").Msgf("%v", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		err := conf.Init(g.log)
		if err != nil {
			g.log.Fatal("golem").Msgf("%v", err)
			os.Exit(1)
		}
	}

	for _, file := range files {
		conf, err := config.NewConfig(file)
		if err != nil {
			g.log.Fatal("golem").Msgf("%v", err)
			os.Exit(1)
		}
		g.conf.Merge(conf)
	}

	g.store = kv.NewStore(g.log)
	g.tpl = template.NewTemplate(g.conf.Servers, *g.conf.Vars, g.store)

	if g.conf.LogLevel != nil {
		g.log = logger.NewCLILogger(int(*g.conf.LogLevel), 12)
	}

	err = g.conf.Recipes.Prepare(g.log, g.store)
	if err != nil {
		g.log.Fatal("golem").Msgf("%v", err)
		os.Exit(1)
	}
	g.Run()

	_ = g.store.Close()
}

func (g *Golem) Run() {
	switch g.cli.Recipe {
	case "":
		g.conf.Recipes.Display(g.log, g.tpl, g.cli.Param1)
	case "version":
		g.log.Highlight("golem").Msgf("version: %s", version)
	case "list":
		g.conf.Recipes.Display(g.log, g.tpl, g.cli.Param1)
	case "servers":
		g.conf.Servers.Display(g.log, g.cli.Param1)
	case "kv":
		switch g.cli.Param1 {
		case "":
			g.store.Display(g.log, "")
		case "list":
			g.store.Display(g.log, "list")
		default:
			switch g.cli.Param2 {
			case "set":
				g.store.SetUserValue(g.cli.Param1)
			case "rand32":
				g.store.SetValue(g.cli.Param1, "rand32")
			case "delete":
				g.store.DeleteValue(g.cli.Param1)
			default:
				g.store.GetValue(g.cli.Param1)
			}
		}
	default:
		g.RunRecipe(g.cli.Recipe)
	}
}

func (g *Golem) RunRecipe(name string) {
	recipe, err := g.conf.Recipes.Search(name)
	if err != nil {
		g.log.Fatal("golem").Msgf("%v", err)
		return
	}

	err = recipe.PrepareForExecution(g.log, g.tpl)
	if err != nil {
		g.log.Fatal("golem").Msgf("%v", err)
		return
	}
	recipe.Display(g.log, g.tpl, "")
	recipe.AskPermission(g.log)
	recipe.Execute(g.log, g.conf.Servers, int(*g.conf.MaxParallelProcesses))
}
