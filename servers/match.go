package servers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/betas-in/utils"
)

type Match struct {
	Attribute string
	Operator  string
	Value     string
}

func NewMatch(attribute, operator, value string) *Match {
	return &Match{
		Attribute: attribute,
		Operator:  operator,
		Value:     value,
	}
}

func (m *Match) Find(c []Server) ([]Server, error) {
	var servers []Server
	for _, s := range c {
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

func (m *Match) server(s Server) (bool, error) {
	if s.Name == "" {
		return false, nil
	}
	switch m.Attribute {
	case "name":
		return m.string("name", s.Name)
	case "public_ip":
		return m.string("public_ip", s.PublicIP)
	case "private_ip":
		return m.string("private_ip", s.PrivateIP)
	case "hostname":
		return m.array("hostname", s.HostName)
	case "user":
		return m.string("user", s.User)
	case "port":
		return m.int("port", s.Port)
	case "tags":
		return m.array("tags", s.Tags)
	default:
		return false, fmt.Errorf("servers does not support attribute %s", m.Attribute)
	}
}

func (m *Match) array(oftype string, list []string) (bool, error) {
	contains := utils.Array().Contains(list, m.Value, true)
	switch m.Operator {
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
	switch m.Operator {
	case "=":
		if name == m.Value {
			return true, nil
		}
	case "!=":
		if name != m.Value {
			return true, nil
		}
	case "like":
		if strings.Contains(name, m.Value) {
			return true, nil
		}
	default:
		// #TODO Support nil match type
		return false, fmt.Errorf("%s only supports ['=', '!=', 'like'] operators", oftype)
	}
	return false, nil
}

func (m *Match) int(oftype string, name int) (bool, error) {
	valueInt, err := strconv.Atoi(m.Value)
	if err != nil {
		return false, nil
	}
	switch m.Operator {
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
