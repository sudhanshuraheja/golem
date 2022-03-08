package config

import (
	"testing"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/servers"
)

func TestServerHCL(t *testing.T) {
	path := "./../testdata/servers.hcl"
	conf, err := NewConfig(path)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 4, len(conf.Servers))
	test1, err := conf.Servers.Search(servers.Match{
		Attribute: "name",
		Operator:  "=",
		Value:     "test4",
	})
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "1.2.3.4", *test1[0].PublicIP)
	utils.Test().Equals(t, "10.11.12.13", *test1[0].PrivateIP)
	utils.Test().Equals(t, 2, len(*test1[0].HostName))
	utils.Test().Equals(t, "sudhanshu", test1[0].User)
	utils.Test().Equals(t, 22, test1[0].Port)
	utils.Test().Equals(t, 4, len(*test1[0].Tags))
}

func TestServerProviderHCL(t *testing.T) {
	path := "./../testdata/serverprovider.hcl"
	conf, err := NewConfig(path)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 3, len(conf.Servers))
	test1, err := conf.Servers.Search(servers.Match{
		Attribute: "name",
		Operator:  "=",
		Value:     "skye",
	})
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "1.2.3.6", *test1[0].PublicIP)
	utils.Test().Equals(t, "10.11.12.15", *test1[0].PrivateIP)
	utils.Test().Equals(t, "root", test1[0].User)
	utils.Test().Equals(t, 22, test1[0].Port)
	utils.Test().Equals(t, 2, len(*test1[0].Tags))
}

func TestRecipesHCL(t *testing.T) {
	path := "./../testdata/recipes.hcl"
	conf, err := NewConfig(path)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 1, len(conf.Servers))
	test1, err := conf.Servers.Search(servers.Match{
		Attribute: "name",
		Operator:  "=",
		Value:     "test1",
	})
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "1.2.3.4", *test1[0].PublicIP)
	utils.Test().Equals(t, "10.11.12.13", *test1[0].PrivateIP)
	utils.Test().Equals(t, "sudhanshu", test1[0].User)
	utils.Test().Equals(t, 22, test1[0].Port)
	utils.Test().Equals(t, 1, len(*test1[0].Tags))

	utils.Test().Equals(t, 2, len(*conf.Vars))
	utils.Test().Equals(t, "golem", (*conf.Vars)["APP"])
}
