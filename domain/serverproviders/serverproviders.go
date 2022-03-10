package serverproviders

import (
	"fmt"

	"github.com/sudhanshuraheja/golem/domain/servers"
)

type ServerProviders []ServerProvider

func (s *ServerProviders) Append(sp ServerProvider) {
	*s = append(*s, sp)
}

func (s *ServerProviders) Merge(sps ServerProviders) {
	*s = append(*s, sps...)
}

func (s *ServerProviders) Parse() (servers.Servers, error) {
	_srvs := servers.Servers{}
	_dips := servers.DomainIPs{}

	if s != nil {
		for _, sp := range *s {
			switch sp.Name {
			case "terraform":
				srvs, dips, err := sp.Parse()
				if err != nil {
					return _srvs, err
				}
				_srvs.Merge(srvs)
				_dips.Merge(dips)
			default:
				return _srvs, fmt.Errorf("server_providers label only supports ['terraform']")
			}
		}
		return _srvs, nil
	}

	_srvs.MergesHosts(&_dips)
	return _srvs, nil
}
