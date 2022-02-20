package kitchen

import (
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/recipes"
)

type Kitchen struct {
	conf *config.Config
}

func NewKitchen(configPath string) *Kitchen {
	conf := config.NewConfig(configPath)
	conf.ResolveServerProvider()
	return &Kitchen{conf: conf}
}

func (k *Kitchen) Exec(recipe string) {
	log.Announcef("%s | running recipe", recipe)
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

		log.Errorf("kitchen | the recipe <%s> was not found, please add it to golem.hcl and try again", recipe)
		log.Infof("Here are the recipes that you can use with 'golem recipe-name'")
		recipes.List(k.conf)
	}
}
