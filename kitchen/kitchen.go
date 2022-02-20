package kitchen

import (
	"os"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/recipes"
)

type Kitchen struct {
	conf *config.Config
}

func NewKitchen(configPath string) *Kitchen {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		if err.Error() == "config file does not exist" {
			recipes.Init()
		} else {
			log.Errorf("%v", err)
		}
		conf, err = config.NewConfig(configPath)
		if err != nil {
			log.Errorf("%v", err)
			os.Exit(1)
		}
	}
	conf.ResolveServerProvider()
	return &Kitchen{conf: conf}
}

func (k *Kitchen) Exec(recipe string) {
	if recipe != "" && k.conf != nil && k.conf.MaxParallelProcesses != nil {
		log.Announcef("%s | running recipe with %d routines", recipe, *k.conf.MaxParallelProcesses)
	}
	switch recipe {
	case "init":
		recipes.Init()
	case "list":
		recipes.List(k.conf)
	case "servers":
		recipes.Servers(k.conf)
	default:
		if recipes.Exists(k.conf, recipe) {
			recipes.Run(k.conf, recipe)
			return
		}

		if recipe != "" {
			log.Errorf("kitchen | the recipe <%s> was not found, please add it to golem.hcl and try again", recipe)
		}
		log.MinorSuccessf("Here are the recipes that you can use with '$ golem recipe-name'\n")
		recipes.List(k.conf)
		log.MinorSuccessf("\nYou can add more recipes to '~/.golem/golem.hcl'")

	}
}
