package recipes

import (
	"strconv"
	"strings"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

type Match struct {
	match config.Match
}

func NewMatch(m config.Match) *Match {
	return &Match{match: m}
}

func (m *Match) Find(c *config.Config) []config.Server {
	var servers []config.Server
	for _, s := range c.Servers {
		if m.server(s) {
			servers = append(servers, s)
		}
	}
	return servers
}

func (m *Match) server(s config.Server) bool {
	if s.Name == "" {
		return false
	}
	switch m.match.Attribute {
	case "name":
		return m.string("name", s.Name)
	case "public_ip":
		if s.PublicIP == nil {
			return false
		}
		return m.string("public_ip", *s.PublicIP)
	case "private_ip":
		if s.PrivateIP == nil {
			return false
		}
		return m.string("private_ip", *s.PrivateIP)
	case "hostname":
		if s.HostName == nil {
			return false
		}
		return m.string("hostname", *s.HostName)
	case "user":
		return m.string("user", s.User)
	case "port":
		return m.int("port", s.Port)
	case "tags":
		return m.array("tags", s.Tags)
	default:
		log.Errorf("servers does not support attribute %s", m.match.Attribute)
	}
	return false
}

func (m *Match) array(oftype string, list []string) bool {
	contains := utils.ArrayContains(list, m.match.Value, true)
	switch m.match.Operator {
	case "contains":
		if contains > -1 {
			return true
		}
	case "not-contains":
		if contains == -1 {
			return true
		}
	default:
		log.Errorf("%s only supports ['contains', 'not-contains'] operators", oftype)
	}
	return false
}

func (m *Match) string(oftype, name string) bool {
	switch m.match.Operator {
	case "=":
		if name == m.match.Value {
			return true
		}
	case "!=":
		if name != m.match.Value {
			return true
		}
	case "like":
		if strings.Contains(name, m.match.Value) {
			return true
		}
	default:
		log.Errorf("%s only supports ['=', '!=', 'like'] operators", oftype)
	}
	return false
}

func (m *Match) int(oftype string, name int) bool {
	valueInt, err := strconv.Atoi(m.match.Value)
	if err != nil {
		return false
	}
	switch m.match.Operator {
	case "=":
		if name == valueInt {
			return true
		}
	case "!=":
		if name != valueInt {
			return true
		}
	case ">":
		if name > valueInt {
			return true
		}
	case ">=":
		if name >= valueInt {
			return true
		}
	case "<":
		if name < valueInt {
			return true
		}
	case "<=":
		if name <= valueInt {
			return true
		}
	default:
		log.Errorf("%s only supports ['=', '!=', '>', '>=', '<', '<='] operators", oftype)
	}
	return false
}
