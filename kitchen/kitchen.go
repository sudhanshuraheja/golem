package kitchen

import (
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
	default:
		recipes.Servers(k.conf)
	}
}
