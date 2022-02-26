package kitchen

import (
	"os"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/recipes"
)

const (
	version = "v0.1.0"
)

type Kitchen struct {
	conf *config.Config
}

func NewKitchen(configPath string) *Kitchen {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		if err.Error() == "config file does not exist" {
			recipes.NewRecipes(nil).Init()
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
		log.Announcef("%s | running recipe with max %d routines", recipe, *k.conf.MaxParallelProcesses)
	}
	r := recipes.NewRecipes(k.conf)
	switch recipe {
	case "":
		log.MinorSuccessf("We found these recipes in '~/.golem/golem.hcl'")
		r.List()
	case "version":
		log.MinorSuccessf("golem version: %s", version)
	case "init":
		r.Init()
	case "list":
		r.List()
	case "servers":
		r.Servers()
	default:
		r.Run(recipe)
	}
}
