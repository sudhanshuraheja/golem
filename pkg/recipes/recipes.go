package recipes

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Context struct {
	Pid string
}

type Config struct {
	Servers ServersConfig `hcl:"servers,block"`
}

type ServersConfig struct {
	Server []Server `hcl:"server,block"`
}

type Server struct {
	Name string `hcl:"name,label"`
	IP   string `hcl:"ip"`
	User string `hcl:"user"`
	Port int    `hcl:"port"`
}

func Start(configPath, recipe string) {
	if configPath == "" {
		configPath = "golem.hcl"
	}
	fmt.Println(configPath, recipe)

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(configPath)
	if diags.HasErrors() {
		log.Fatalf("parse error: %v", diags)
	}

	// ctx := &hcl.EvalContext{
	// 	Variables: map[string]cty.Value{
	// 		"id":   cty.StringVal("Emintrude"),
	// 		"port": cty.NumberIntVal(22),
	// 	},
	// }

	var c Config
	diags = gohcl.DecodeBody(f.Body, nil, &c)
	if diags.HasErrors() {
		log.Fatalf("parse body error: %v", diags)
	}

	fmt.Printf("%+v", c)
}
