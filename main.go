package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sudhanshuraheja/golem/kitchen"
)

func main() {
	conf := kitchen.Config{}
	arg.MustParse(&conf)
	kitchen.NewKitchen(&conf)
}
