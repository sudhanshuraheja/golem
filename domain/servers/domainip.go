package servers

type DomainIP struct {
	Host string
	IP   string
}

func NewDomainIP(host, ip string) *DomainIP {
	return &DomainIP{host, ip}
}
