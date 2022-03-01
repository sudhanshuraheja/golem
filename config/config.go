package config

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Config struct {
	ServerProviders      []ServerProvider   `hcl:"server_provider,block"`
	Servers              []Server           `hcl:"server,block"`
	Recipes              []Recipe           `hcl:"recipe,block"`
	LogLevel             *int               `hcl:"loglevel"`
	MaxParallelProcesses *int               `hcl:"max_parallel_processes"`
	Vars                 *map[string]string `hcl:"vars"`
}

type ServerProvider struct {
	Name   string   `hcl:"name,label"`
	Config []string `hcl:"config"`
	User   string   `hcl:"user"`
	Port   int      `hcl:"port"`
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
	Type      string     `hcl:"type,label"`
	Match     *Match     `hcl:"match,block"`
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

func NewConfig(path string) (*Config, error) {
	var conf Config

	parser := hclparse.NewParser()

	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, diags
	}

	diags = gohcl.DecodeBody(f.Body, nil, &conf)
	if diags.HasErrors() {
		return nil, diags
	}

	if conf.MaxParallelProcesses == nil {
		maxParallelProcs := 4
		conf.MaxParallelProcesses = &maxParallelProcs
	}

	return &conf, nil
}

func (c *Config) ResolveServerProvider() error {
	for _, sp := range c.ServerProviders {
		switch sp.Name {
		case "terraform":
			for _, cf := range sp.Config {
				spt := ServerProviderTerraform{}
				srvs, iph, err := spt.GetServers(cf, sp.User, sp.Port)
				if err != nil {
					return fmt.Errorf("could not load servers from tfstate %s: %v", cf, err)
				}
				c.Servers = append(c.Servers, srvs...)
				mergeIPHostnames(&c.Servers, iph)
			}
		default:
			return fmt.Errorf("server_providers label only supports ['terraform']")
		}
	}
	return nil
}

func mergeIPHostnames(servers *[]Server, iph IPHostNames) {
	for ip, hostnames := range iph {
		for i, server := range *servers {
			if server.PublicIP != nil && *server.PublicIP == ip {
				hn := strings.Join(hostnames, ", ")
				(*servers)[i].HostName = &hn
				delete(iph, ip)
			}
		}
	}

	for ip, hostnames := range iph {
		srv := Server{}
		ipToUse := ip
		srv.PublicIP = &ipToUse
		hn := strings.Join(hostnames, ", ")
		srv.HostName = &hn
		*servers = append(*servers, srv)
	}
}
