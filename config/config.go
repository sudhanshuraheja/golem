package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sudhanshuraheja/golem/pkg/log"
)

type Config struct {
	ServerProviders      []ServerProvider `hcl:"server_provider,block"`
	Servers              []Server         `hcl:"server,block"`
	Recipe               []Recipe         `hcl:"recipe,block"`
	LogLevel             *string          `hcl:"loglevel"`
	MaxParallelProcesses *int             `hcl:"max_parallel_processes"`
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

func getConfFilePath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find user's home directory: %v", err)
	}
	return fmt.Sprintf("%s/.golem/golem.hcl", dirname), nil
}

func NewConfig(configPath string) *Config {
	var conf Config

	if configPath == "" {
		var err error
		configPath, err = getConfFilePath()
		if err != nil {
			log.Errorf("%v", err)
			return nil
		}
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

	if conf.MaxParallelProcesses == nil {
		maxParallelProcs := 4
		conf.MaxParallelProcesses = &maxParallelProcs
	}

	return &conf
}

func (c *Config) ResolveServerProvider() {
	startTime := time.Now()
	for _, sp := range c.ServerProviders {
		switch sp.Name {
		case "terraform":
			for _, cf := range sp.Config {
				spt := ServerProviderTerraform{}
				srvs, iph, err := spt.GetServers(cf, sp.User, sp.Port)
				if err != nil {
					log.Errorf("config | could not load servers from tfstate %s: %v", cf, err)
					continue
				}
				log.Debugf("config | found %d servers in %s", len(srvs), cf)
				c.Servers = append(c.Servers, srvs...)
				mergeIPHostnames(&c.Servers, iph)
			}
		default:
			log.Errorf("config | server_providers label only supports ['terraform']")
		}
	}
	log.Debugf("resolved service provider in %s", time.Since(startTime))
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
