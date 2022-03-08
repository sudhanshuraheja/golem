package servers

import (
	"fmt"
	"strings"

	"github.com/betas-in/logger"
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

func (s *Server) Log(log *logger.CLILogger, query string) {
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

func Cyan(name, value string) string {
	return fmt.Sprintf("%s %s ", logger.Cyan(name), value)
}
