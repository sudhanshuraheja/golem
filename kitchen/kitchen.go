package kitchen

import (
	"github.com/fatih/color"
	"github.com/sudhanshuraheja/golem/config"
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
	case "tflist":
		recipes.TerraformResources(k.conf, "")
	case "tflistall":
		recipes.TerraformResources(k.conf, "all")
	case "servers":
		recipes.Servers(k.conf)
	default:
		red := color.New(color.FgRed)
		red.Printf("The recipe %s was not found, using golem servers\n", recipe)
		recipes.Servers(k.conf)
	}
}
