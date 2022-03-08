package servers

type DomainIPs []DomainIP

func (d *DomainIPs) Append(dip DomainIP) {
	if d != nil {
		*d = append(*d, dip)
	}
}

func (d *DomainIPs) Merge(dips DomainIPs) {
	if d != nil {
		*d = append(*d, dips...)
	}
}
