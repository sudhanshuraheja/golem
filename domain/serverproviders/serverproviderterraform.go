package serverproviders

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/servers"
)

type ServerProviderTerraform struct {
	Resources []TFResource
}

type TFResource struct {
	Type      string
	Instances []TFInstance
}

type TFInstance struct {
	Attributes TFAttributes
}

type TFAttributes struct {
	CreatedAt          string             `json:"created_at"`
	IPV4Address        string             `json:"ipv4_address"`
	IPV4AddressPrivate string             `json:"ipv4_address_private"`
	Name               string             `json:"name"`
	Region             string             `json:"region"`
	Tags               []string           `json:"tags"`
	Type               string             `json:"type"`
	FQDN               string             `json:"fqdn"`
	Value              string             `json:"value"`
	InboundRule        []TFDOInboundRule  `json:"inbound_rule"`
	OutboundRule       []TFDOOutboundRule `json:"outbound_rule"`
}

type TFDOInboundRule struct {
	PortRange       string   `json:"port_range"`
	Protocol        string   `json:"protocol"`
	SourceAddresses []string `json:"source_addresses"`
}

type TFDOOutboundRule struct {
	PortRange            string   `json:"port_range"`
	Protocol             string   `json:"protocol"`
	DestinationAddresses []string `json:"destination_addresses"`
}

func (s *ServerProviderTerraform) GetServers(file, user string, port int) (servers.Servers, error) {
	srvs := servers.Servers{}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return srvs, fmt.Errorf("unable to read file: %s: %v", file, err)
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return srvs, fmt.Errorf("unable to unmarshall: %s: %v", file, err)
	}

	for _, tfResource := range s.Resources {
		if utils.Array().Contains([]string{
			"digitalocean_droplet",
		}, tfResource.Type, true) >= 0 {
			for _, tfResourceInstance := range tfResource.Instances {
				srv := servers.Server{}
				srv.Name = tfResourceInstance.Attributes.Name
				srv.PublicIP = &tfResourceInstance.Attributes.IPV4Address
				srv.PrivateIP = &tfResourceInstance.Attributes.IPV4AddressPrivate
				srv.HostName = nil
				srv.User = user
				srv.Port = port
				srv.Tags = tfResourceInstance.Attributes.Tags
				srvs.Append(srv)
			}
		}
	}

	return srvs, nil
}

func (s *ServerProviderTerraform) GetDomainIP(file string) (servers.DomainIPs, error) {
	domainIPs := servers.DomainIPs{}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return domainIPs, fmt.Errorf("unable to read file: %s: %v", file, err)
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return domainIPs, fmt.Errorf("unable to unmarshall: %s: %v", file, err)
	}

	for _, tfResource := range s.Resources {
		if utils.Array().Contains([]string{
			"digitalocean_record",
		}, tfResource.Type, true) >= 0 {
			for _, tfResourceInstance := range tfResource.Instances {
				if tfResourceInstance.Attributes.Type == "A" {
					dip := servers.DomainIP{}
					dip.Host = tfResourceInstance.Attributes.FQDN
					dip.IP = tfResourceInstance.Attributes.Value
					domainIPs.Append(dip)
				}
			}
		}
	}

	return domainIPs, nil
}
