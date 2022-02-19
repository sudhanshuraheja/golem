package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sudhanshuraheja/golem/kitchen"
)

var args struct {
	Recipe string `arg:"positional"`
	Config string `arg:"-c,--conf" help:"config folder, can be a file ./golem.hcl or folder ./recipes/"`
}

func main() {
	arg.MustParse(&args)
	kitchen := kitchen.NewKitchen(args.Config)
	kitchen.Exec(args.Recipe)
}
