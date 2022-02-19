package recipes

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/sudhanshuraheja/golem/config"
)

type TFRoot struct {
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

func TerraformResources(c *config.Config, filter string) {
	tb := NewTable("Type", "Name", "IPV4Address", "IPV4AddressPrivate", "Region", "Tags", "Type", "Value", "CreatedAt")

	for _, t := range *c.Terraform {

		bytes, err := os.ReadFile(t)
		if err != nil {
			log.Fatalf("unable to read file: %s: %v", t, err)
		}

		var tf TFRoot
		err = json.Unmarshal(bytes, &tf)
		if err != nil {
			log.Fatalf("unable to unmarshall: %s: %v", t, err)
		}

		for _, tfr := range tf.Resources {
			resourceType := tfr.Type
			resourceType = strings.Replace(resourceType, "digitalocean", "do", -1)
			resourceType = strings.Replace(resourceType, "kubernetes", "k8s", -1)

			for _, tfi := range tfr.Instances {
				name := tfi.Attributes.Name
				if resourceType == "do_record" {
					name = tfi.Attributes.FQDN
				}

				value := tfi.Attributes.Value
				if len(value) > 20 {
					value = value[:20]
				}

				createdAt := tfi.Attributes.CreatedAt
				if len(createdAt) > 10 {
					createdAt = createdAt[:10]
				}

				if filter == "" {
					if strings.Contains(resourceType, "_domain") || strings.Contains(resourceType, "_record") {
						continue
					}
				}

				tb.Row(
					resourceType,
					name,
					tfi.Attributes.IPV4Address,
					tfi.Attributes.IPV4AddressPrivate,
					tfi.Attributes.Region,
					strings.Join(tfi.Attributes.Tags, ","),
					tfi.Attributes.Type,
					value,
					createdAt,
				)
			}
		}
	}

	tb.Display()
}
