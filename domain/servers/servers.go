package servers

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
					(*s)[i].HostName = append((*s)[i].HostName, dip.Host)
					// found = true
				}
			}
			if !found {
				srv := Server{}
				srv.PublicIP = &dip.IP
				srv.HostName = []string{dip.Host}
				*s = append(*s, srv)
			}
		}
	}
}
