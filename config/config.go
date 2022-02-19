package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sudhanshuraheja/golem/pkg/log"
)

type Config struct {
	Servers   ServersConfig `hcl:"servers,block"`
	Recipe    []Recipe      `hcl:"recipe,block"`
	Terraform *[]string     `hcl:"terraform"`
	LogLevel  *string       `hcl:"loglevel"`
}

type ServersConfig struct {
	Server []Server `hcl:"server,block"`
}

type Server struct {
	Name      string   `hcl:"name,label"`
	PublicIP  *string  `hcl:"public_ip"`
	PrivateIP *string  `hcl:"private_ip"`
	HostName  *string  `hcl:"hostname"`
	User      string   `hcl:"user"`
	Port      int      `hcl:"port"`
	Tags      []string `hcl:"tags"`
}

type Recipe struct {
	Name   string   `hcl:"name,label"`
	Target []string `hcl:"target"`
}

func NewConfig(configPath string) *Config {
	var conf Config

	if configPath == "" {
		configPath = "golem.hcl"
	}

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(configPath)
	if diags.HasErrors() {
		fmt.Printf("parse error: %v", diags)
		os.Exit(1)
	}

	diags = gohcl.DecodeBody(f.Body, nil, &conf)
	if diags.HasErrors() {
		fmt.Printf("parse body error: %v", diags)
		os.Exit(1)
	}

	log.SetLogLevel(conf.LogLevel)
	return &conf
}
