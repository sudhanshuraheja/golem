package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sudhanshuraheja/golem/pkg/recipes"
)

var args struct {
	Recipe string `arg:"positional"`
	Config string `arg:"-c,--conf" help:"config folder, can be a file ./golem.hcl or folder ./recipes/"`
}

func main() {
	arg.MustParse(&args)
	recipes.Start(args.Config, args.Recipe)
}
