package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sudhanshuraheja/golem/domain/natives"
	"github.com/sudhanshuraheja/golem/domain/recipes"
	"github.com/sudhanshuraheja/golem/domain/serverproviders"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/vars"
)

type Config struct {
	Servers              servers.Servers                 `hcl:"server,block"`
	ServerProviders      serverproviders.ServerProviders `hcl:"server_provider,block"`
	Recipes              recipes.Recipes                 `hcl:"recipe,block"`
	LogLevel             *natives.Int                    `hcl:"loglevel"`
	MaxParallelProcesses *natives.Int                    `hcl:"max_parallel_processes"`
	Vars                 *vars.Vars                      `hcl:"vars"`
}

func NewConfig(path string) *Config {
	var conf Config

	parser := hclparse.NewParser()

	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		showHCLDiagnostics(parser, diags)
		return nil
	}

	diags = gohcl.DecodeBody(f.Body, nil, &conf)
	if diags.HasErrors() {
		showHCLDiagnostics(parser, diags)
		return nil
	}

	return &conf
}

func showHCLDiagnostics(parser *hclparse.Parser, diags hcl.Diagnostics) {
	wr := hcl.NewDiagnosticTextWriter(
		os.Stdout,
		parser.Files(),
		80,
		true,
	)

	for _, diag := range diags {
		err := wr.WriteDiagnostic(diag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(1)
}

func (c *Config) ParseServerProviders() error {
	srvs, err := c.ServerProviders.Parse()
	if err != nil {
		return err
	}
	c.Servers.Merge(srvs)
	return nil
}
