package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sudhanshuraheja/golem/pkg/log"
)

type Config struct {
	Servers   []Server  `hcl:"server,block"`
	Recipe    []Recipe  `hcl:"recipe,block"`
	Terraform *[]string `hcl:"terraform"`
	LogLevel  *string   `hcl:"loglevel"`
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
	Name      string     `hcl:"name,label"`
	Type      string     `hcl:"type"`
	Match     Match      `hcl:"match,block"`
	Artifacts []Artifact `hcl:"artifact,block"`
	Commands  []string   `hcl:"commands"`
}

type Match struct {
	Attribute string `hcl:"attribute"`
	Operator  string `hcl:"operator"`
	Value     string `hcl:"value"`
}

type Artifact struct {
	Source      string `hcl:"source"`
	Destination string `hcl:"destination"`
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
