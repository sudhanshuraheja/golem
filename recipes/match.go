package recipes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

type Match struct {
	match config.Match
}

func NewMatch(m config.Match) *Match {
	return &Match{match: m}
}

func (m *Match) Find(c *config.Config) ([]config.Server, error) {
	var servers []config.Server
	for _, s := range c.Servers {
		matched, err := m.server(s)
		if err != nil {
			return servers, err
		}
		if matched {
			servers = append(servers, s)
		}
	}
	return servers, nil
}

func (m *Match) server(s config.Server) (bool, error) {
	if s.Name == "" {
		return false, nil
	}
	switch m.match.Attribute {
	case "name":
		return m.string("name", s.Name)
	case "public_ip":
		if s.PublicIP == nil {
			return false, nil
		}
		return m.string("public_ip", *s.PublicIP)
	case "private_ip":
		if s.PrivateIP == nil {
			return false, nil
		}
		return m.string("private_ip", *s.PrivateIP)
	case "hostname":
		if s.HostName == nil {
			return false, nil
		}
		return m.string("hostname", *s.HostName)
	case "user":
		return m.string("user", s.User)
	case "port":
		return m.int("port", s.Port)
	case "tags":
		return m.array("tags", s.Tags)
	default:
		return false, fmt.Errorf("servers does not support attribute %s", m.match.Attribute)
	}
}

func (m *Match) array(oftype string, list []string) (bool, error) {
	contains := utils.Array().Contains(list, m.match.Value, true)
	switch m.match.Operator {
	case "contains":
		if contains > -1 {
			return true, nil
		}
	case "not-contains":
		if contains == -1 {
			return true, nil
		}
	default:
		return false, fmt.Errorf("%s only supports ['contains', 'not-contains'] operators", oftype)
	}
	return false, nil
}

func (m *Match) string(oftype, name string) (bool, error) {
	switch m.match.Operator {
	case "=":
		if name == m.match.Value {
			return true, nil
		}
	case "!=":
		if name != m.match.Value {
			return true, nil
		}
	case "like":
		if strings.Contains(name, m.match.Value) {
			return true, nil
		}
	default:
		return false, fmt.Errorf("%s only supports ['=', '!=', 'like'] operators", oftype)
	}
	return false, nil
}

func (m *Match) int(oftype string, name int) (bool, error) {
	valueInt, err := strconv.Atoi(m.match.Value)
	if err != nil {
		return false, nil
	}
	switch m.match.Operator {
	case "=":
		if name == valueInt {
			return true, nil
		}
	case "!=":
		if name != valueInt {
			return true, nil
		}
	case ">":
		if name > valueInt {
			return true, nil
		}
	case ">=":
		if name >= valueInt {
			return true, nil
		}
	case "<":
		if name < valueInt {
			return true, nil
		}
	case "<=":
		if name <= valueInt {
			return true, nil
		}
	default:
		return false, fmt.Errorf("%s only supports ['=', '!=', '>', '>=', '<', '<='] operators", oftype)
	}
	return false, nil
}
