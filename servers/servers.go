package servers

import (
	"fmt"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
)

var (
	tiny = 50
)

type Server struct {
	Name      string
	PublicIP  string
	PrivateIP string
	HostName  []string
	User      string
	Port      int
	Tags      []string
}

func NewServers(conf []config.Server) []Server {
	servers := []Server{}
	for _, s := range conf {
		server := Server{}

		server.Name = s.Name

		if s.PublicIP != nil {
			server.PublicIP = *s.PublicIP
		}

		if server.Name == "" {
			server.Name = server.PublicIP
		}

		if s.PrivateIP != nil {
			server.PrivateIP = *s.PrivateIP
		}

		server.HostName = s.HostName
		server.User = s.User
		if server.User == "" {
			server.User = "root"
		}

		server.Port = s.Port
		if server.Port == 0 {
			server.Port = 22
		}

		server.Tags = s.Tags

		servers = append(servers, server)
	}
	return servers
}

func (s *Server) Display(log *logger.CLILogger, query string) {

	if len(query) > 0 {
		if !strings.Contains(s.Name, query) {
			return
		}
	}

	publicIP := s.PublicIP
	if publicIP != "" {
		publicIP = fmt.Sprintf("%s %s ", logger.Cyan("publicIP"), publicIP)
	}

	privateIP := s.PrivateIP
	if privateIP != "" {
		privateIP = fmt.Sprintf("%s %s ", logger.Cyan("privateIP"), privateIP)
	}

	hostnames := strings.Join(s.HostName, ", ")
	tags := strings.Join(s.Tags, ", ")

	log.Info(s.Name).Msgf(
		"%s%s%s%s",
		fmt.Sprintf("%s %s ", logger.Cyan("user"), s.User),
		fmt.Sprintf("%s %d ", logger.Cyan("port"), s.Port),
		publicIP,
		privateIP,
	)

	if hostnames != "" {
		log.Info("").Msgf("%s %s", logger.Cyan("hosts"), hostnames)
	}

	if tags != "" {
		log.Info("").Msgf("%s %s", logger.Cyan("tags"), tags)
	}

}
