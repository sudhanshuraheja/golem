package servers

import (
	"fmt"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Server struct {
	Name      string   `hcl:"name,label"`
	PublicIP  *string  `hcl:"public_ip"`
	PrivateIP *string  `hcl:"private_ip"`
	HostName  []string `hcl:"hostname"`
	User      string   `hcl:"user"`
	Port      int      `hcl:"port"`
	Tags      []string `hcl:"tags"`
}

func (s *Server) Display(log *logger.CLILogger, query string) {
	if len(query) > 0 {
		if !strings.Contains(s.Name, query) {
			return
		}
	}

	publicIP := ""
	if s.PublicIP != nil {
		publicIP = Cyan("publicIP", *s.PublicIP)
	}

	privateIP := ""
	if s.PrivateIP != nil {
		privateIP = Cyan("privateIP", *s.PrivateIP)
	}

	log.Info(s.Name).Msgf("%s%s%s%s", Cyan("user", s.User), fmt.Sprintf("%s %d ", logger.Cyan("port"), s.Port), publicIP, privateIP)

	hostnames := strings.Join(s.HostName, ", ")
	if hostnames != "" {
		log.Info("").Msgf("%s %s", logger.Cyan("hosts"), hostnames)
	}

	tags := strings.Join(s.Tags, ", ")
	if tags != "" {
		log.Info("").Msgf("%s %s", logger.Cyan("tags"), tags)
	}

}

func (s *Server) GetHostName() (string, error) {
	host := ""

	switch {
	case s.PublicIP != nil:
		host = *s.PublicIP
	case len(s.HostName) > 0:
		host = s.HostName[0]
	default:
		return host, fmt.Errorf("could not find a public ip or a hostname in config")
	}

	return host, nil
}

func (s *Server) Search(m Match) (bool, error) {
	if s.Name == "" {
		return false, nil
	}

	switch m.Attribute {
	case "name":
		return m.CompareString(s.Name)
	case "public_ip":
		return m.CompareString(localutils.StringPtrValue(s.PublicIP, ""))
	case "private_ip":
		return m.CompareString(localutils.StringPtrValue(s.PrivateIP, ""))
	case "hostname":
		return m.CompareStringArray(s.HostName)
	case "user":
		return m.CompareString(s.User)
	case "port":
		return m.CompareInt(s.Port)
	case "tags":
		return m.CompareStringArray(s.Tags)
	default:
		return false, fmt.Errorf("servers does not support attribute %s", m.Attribute)
	}
}

func Cyan(name, value string) string {
	return fmt.Sprintf("%s %s ", logger.Cyan(name), value)
}
