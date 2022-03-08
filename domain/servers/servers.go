package servers

import "github.com/betas-in/logger"

type Servers []Server

func (s *Servers) Append(srv Server) {
	if s != nil {
		*s = append(*s, srv)
	}
}

func (s *Servers) Merge(srvs Servers) {
	if s != nil {
		*s = append(*s, srvs...)
	}
}

func (s *Servers) MergesHosts(dips *DomainIPs) {
	if dips != nil {
		for _, dip := range *dips {
			found := false
			for i, srv := range *s {
				if srv.PublicIP != nil && *srv.PublicIP == dip.IP {
					*((*s)[i].HostName) = append(*((*s)[i].HostName), dip.Host)
					// found = true
				}
			}
			if !found {
				srv := Server{}
				srv.PublicIP = &dip.IP
				srv.HostName = &[]string{dip.Host}
				*s = append(*s, srv)
			}
		}
	}
}

func (s *Servers) Search(m Match) (Servers, error) {
	servers := Servers{}
	for _, srv := range *s {
		matched, err := srv.Search(m)
		if err != nil {
			return servers, err
		}
		if matched {
			servers.Append(srv)
		}
	}
	return servers, nil
}

func (s *Servers) Display(log *logger.CLILogger, query string) {
	log.Announce("").Msgf("list of all connected servers")
	for _, srv := range *s {
		srv.Display(log, query)
	}
}
