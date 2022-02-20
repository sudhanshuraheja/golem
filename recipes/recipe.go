package recipes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func List(c *config.Config) {
	tb := log.NewTable("Name", "Match", "Artifacts", "Commands")
	for _, r := range c.Recipe {
		tb.Row(
			r.Name,
			fmt.Sprintf("%s %s %s", r.Match.Attribute, r.Match.Operator, r.Match.Value),
			len(r.Artifacts),
			len(r.Commands),
		)
	}
	// Add system defined
	tb.Row("servers", "local only", 0, 0)
	tb.Row("tflist", "local only", 0, 0)
	tb.Row("tflistall", "local only", 0, 0)
	tb.Display()
}

func Exists(c *config.Config, name string) bool {
	for _, r := range c.Recipe {
		if r.Name == name {
			return true
		}
	}
	return false
}

func Run(c *config.Config, name string) {
	var recipe config.Recipe
	for i, r := range c.Recipe {
		if r.Name == name {
			recipe = c.Recipe[i]
		}
	}

	servers := findMatchingServers(c, recipe.Match)
	serverNames := []string{}
	for _, s := range servers {
		serverNames = append(serverNames, s.Name)
	}
	log.Announcef("%s | found %d matching servers - %s", recipe.Name, len(servers), strings.Join(serverNames, ", "))

	switch recipe.Type {
	case "exec":
		for _, s := range servers {
			SSHRun(&s, recipe.Commands)
		}
	default:
		log.Errorf("recipe only supports ['exec'] types")
	}
}

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
