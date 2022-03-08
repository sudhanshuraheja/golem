package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sudhanshuraheja/golem/golem"
)

func main() {
	conf := golem.Config{}
	arg.MustParse(&conf)
	golem.NewGolem(&conf)
}
