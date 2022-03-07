package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sudhanshuraheja/golem/kitchen"
)

func main() {
	conf := kitchen.CLIConfig{}
	arg.MustParse(&conf)
	kitchen.NewKitchen(&conf)
}
