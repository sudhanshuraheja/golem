package recipes

import (
	"strconv"
	"strings"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func findMatchingServers(c *config.Config, match config.Match) []config.Server {
	var servers []config.Server
	for _, s := range c.Servers {
		if matchServer(s, match) {
			servers = append(servers, s)
		}
	}
	return servers
}

func matchServer(s config.Server, m config.Match) bool {
	if s.Name == "" {
		return false
	}
	switch m.Attribute {
	case "name":
		return matchString("name", s.Name, m.Value, m.Operator)
	case "public_ip":
		if s.PublicIP == nil {
			return false
		}
		return matchString("public_ip", *s.PublicIP, m.Value, m.Operator)
	case "private_ip":
		if s.PrivateIP == nil {
			return false
		}
		return matchString("private_ip", *s.PrivateIP, m.Value, m.Operator)
	case "hostname":
		if s.HostName == nil {
			return false
		}
		return matchString("hostname", *s.HostName, m.Value, m.Operator)
	case "user":
		return matchString("user", s.User, m.Value, m.Operator)
	case "port":
		return matchInt("user", s.Port, m.Value, m.Operator)
	case "tags":
		return matchArray("tags", s.Tags, m.Value, m.Operator)
	default:
		log.Errorf("servers does not support attribute %s", m.Attribute)
	}
	return false
}

func matchArray(oftype string, list []string, value, operator string) bool {
	contains := utils.ArrayContains(list, value, true)
	switch operator {
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

func matchString(oftype, name, value, operator string) bool {
	switch operator {
	case "=":
		if name == value {
			return true
		}
	case "!=":
		if name != value {
			return true
		}
	case "like":
		if strings.Contains(name, value) {
			return true
		}
	default:
		log.Errorf("%s only supports ['=', '!=', 'like'] operators", oftype)
	}
	return false
}

func matchInt(oftype string, name int, value, operator string) bool {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	switch operator {
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
