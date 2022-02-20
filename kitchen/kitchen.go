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
	return &Kitchen{
		conf: config.NewConfig(configPath),
	}
}

func (k *Kitchen) Exec(recipe string) {
	switch recipe {
	case "list":
		recipes.List(k.conf)
	// case "apt-update":
	// 	recipes.AptUpdate(k.conf)
	case "tflist":
		recipes.TerraformResources(k.conf, "")
	case "tflistall":
		recipes.TerraformResources(k.conf, "all")
	case "servers":
		recipes.Servers(k.conf)
	default:
		log.Errorf("kitchen | the recipe <%s> was not found, please add it to golem.hcl and try again", recipe)
	}
}
