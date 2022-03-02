package recipes

import (
	"testing"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

func TestMatcher(t *testing.T) {
	publicIP := "127.0.0.1"
	servers := []config.Server{
		{
			Name:     "one",
			HostName: []string{"one"},
			Port:     22,
			Tags:     []string{"one"},
		},
		{
			Name:     "two",
			PublicIP: &publicIP,
			HostName: []string{"two"},
			Port:     22,
			Tags:     []string{"one", "two"},
		},
		{
			Name:     "three",
			HostName: []string{"three"},
			Port:     22,
			Tags:     []string{"one", "two", "three"},
		},
	}

	match := config.Match{}

	match.Attribute = "tags"
	match.Operator = "contains"
	match.Value = "one"
	found, err := NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 3, len(found))

	match.Value = "two"
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 2, len(found))

	match.Value = "three"
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 1, len(found))

	match.Operator = "not-contains"
	match.Value = "two"
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 1, len(found))

	match.Attribute = "public_ip"
	match.Operator = "="
	match.Value = publicIP
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 1, len(found))

	match.Operator = "!="
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 0, len(found))

	match.Attribute = "port"
	match.Operator = "="
	match.Value = "22"
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 3, len(found))

	match.Operator = "!="
	found, err = NewMatch(match).Find(servers)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 0, len(found))

}
