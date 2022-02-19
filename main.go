package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

var args struct {
	Server    string `arg:"-s,--server" help:"sudhanshu@1.1.1.1"`
	PublicKey string `arg:"-p,--publickey" help:"publickey"`
}

func main() {
	arg.MustParse(&args)
	fmt.Println(args)
}
