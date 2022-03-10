package config

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sudhanshuraheja/golem/domain/recipes"
	"github.com/sudhanshuraheja/golem/domain/serverproviders"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/vars"
)

type Config struct {
	Servers              servers.Servers                 `hcl:"server,block"`
	ServerProviders      serverproviders.ServerProviders `hcl:"server_provider,block"`
	Recipes              recipes.Recipes                 `hcl:"recipe,block"`
	LogLevel             *int                            `hcl:"loglevel"`
	MaxParallelProcesses *int                            `hcl:"max_parallel_processes"`
	Vars                 *vars.Vars                      `hcl:"vars"`
}

func NewConfig(path string) (*Config, error) {
	var conf Config

	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		showHCLDiagnostics(parser, diags)
		return nil, diags
	}

	diags = gohcl.DecodeBody(file.Body, nil, &conf)
	if diags.HasErrors() {
		showHCLDiagnostics(parser, diags)
		return nil, diags
	}

	err := conf.ParseServerProviders()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func showHCLDiagnostics(parser *hclparse.Parser, diags hcl.Diagnostics) {
	writer := hcl.NewDiagnosticTextWriter(
		os.Stdout,
		parser.Files(),
		80,
		true,
	)

	for _, diag := range diags {
		_ = writer.WriteDiagnostic(diag)
	}
}

func (c *Config) ParseServerProviders() error {
	srvs, err := c.ServerProviders.Parse()
	if err != nil {
		return err
	}
	c.Servers.Merge(srvs)
	return nil
}

func (c *Config) Merge(conf *Config) {
	c.Servers.Merge(conf.Servers)
	c.Recipes.Merge(conf.Recipes)

	if conf.LogLevel != nil {
		c.LogLevel = conf.LogLevel
	}

	if conf.MaxParallelProcesses != nil {
		c.MaxParallelProcesses = conf.MaxParallelProcesses
	}

	if conf.Vars != nil {
		if c.Vars == nil {
			c.Vars = vars.NewVars()
		}
		for key, value := range *conf.Vars {
			(*c.Vars)[key] = value
		}
	}
}
