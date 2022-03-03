package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sudhanshuraheja/golem/kitchen"
)

var args struct {
	Recipe string `arg:"positional"`
	Param1 string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)
	kitchen := kitchen.NewKitchen()
	kitchen.Exec(args.Recipe, args.Param1)
}
