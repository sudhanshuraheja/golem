package serverproviders

import (
	"fmt"

	"github.com/sudhanshuraheja/golem/domain/servers"
)

type ServerProvider struct {
	Name   string   `hcl:"name,label"`
	Config []string `hcl:"config"`
	User   string   `hcl:"user"`
	Port   int      `hcl:"port"`
}

func (s *ServerProvider) Parse() (servers.Servers, servers.DomainIPs, error) {
	_srvs := servers.Servers{}
	_dips := servers.DomainIPs{}

	switch s.Name {
	case "terraform":
		for _, cf := range s.Config {
			spt := ServerProviderTerraform{}

			srvs, err := spt.GetServers(cf, s.User, s.Port)
			if err != nil {
				return _srvs, _dips, fmt.Errorf("could not load servers from tfstate %s: %v", cf, err)
			}

			dips, err := spt.GetDomainIP(cf)
			if err != nil {
				return _srvs, _dips, fmt.Errorf("could not load ips from tfstate %s: %v", cf, err)
			}

			_srvs.Merge(srvs)
			_dips.Merge(dips)
		}
		return _srvs, _dips, nil
	default:
		return _srvs, _dips, fmt.Errorf("server_providers label only supports ['terraform']")
	}
}
