package recipes

import (
	"fmt"
	"os"
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

func Init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("init | could not find user's home directory: %v", err)
		return
	}
	confDir := fmt.Sprintf("%s/.golem", dirname)

	err = os.MkdirAll(confDir, os.FileMode(0755))
	if err != nil {
		log.Errorf("init | could not create conf dir %s: %v", confDir, err)
		return
	}

	confFile := fmt.Sprintf("%s/golem.hcl", confDir)
	_, err = os.Stat(confFile)
	if os.IsNotExist(err) {
		file, err := os.Create(confFile)
		if err != nil {
			log.Errorf("init | error creating conf file %s: %v", confFile, err)
			return
		}
		defer file.Close()
		log.MinorSuccessf("init | conf file created at %s", confFile)
	} else if err != nil {
		log.Errorf("init | error checking conf file %s: %v", confFile, err)
	}
	log.MinorSuccessf("init | conf file already exists at %s", confFile)
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

	for _, a := range recipe.Artifacts {
		log.Infof("→ %s → %s", a.Source, a.Destination)
	}

	for _, c := range recipe.Commands {
		log.Infof("→ $ %s", c)
	}

	answer := log.Questionf("Are you sure you want to continue [y/n]?")
	if utils.ArrayContains([]string{"y", "yes"}, answer, false) == -1 {
		log.Errorf("Quitting, because you said %s", answer)
		os.Exit(0)
	}

	switch recipe.Type {
	case "exec":
		for _, s := range servers {
			ss := SSH{}
			err := ss.Connect(&s)
			if err != nil {
				log.Errorf("%s | %v", s.Name, err)
				continue
			}
			ss.Upload(recipe.Artifacts)
			ss.Run(recipe.Commands)
			ss.Close()
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
