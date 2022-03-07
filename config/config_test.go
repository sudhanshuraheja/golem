package config

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestConfig(t *testing.T) {
	path := "./../testdata/sample.hcl"
	conf := NewConfig(path)

	utils.Test().Equals(t, 1, len(conf.ServerProviders))
	utils.Test().Equals(t, "terraform", conf.ServerProviders[0].Name)
	utils.Test().Equals(t, 0, len(conf.ServerProviders[0].Config))
	utils.Test().Equals(t, "sudhanshu", conf.ServerProviders[0].User)
	utils.Test().Equals(t, 22, conf.ServerProviders[0].Port)

	utils.Test().Equals(t, 1, len(conf.Servers))
	utils.Test().Equals(t, "test-server", conf.Servers[0].Name)
	utils.Test().Equals(t, "127.0.0.1", *conf.Servers[0].PublicIP)
	utils.Test().Equals(t, "127.0.0.1", *conf.Servers[0].PrivateIP)
	utils.Test().Equals(t, 1, len(conf.Servers[0].HostName))
	utils.Test().Equals(t, "sudhanshu", conf.Servers[0].User)
	utils.Test().Equals(t, 22, conf.Servers[0].Port)
	utils.Test().Equals(t, 1, len(conf.Servers[0].Tags))

	utils.Test().Equals(t, 3, *conf.LogLevel)
	utils.Test().Equals(t, 5, *conf.MaxParallelProcesses)

	utils.Test().Equals(t, 1, len(*conf.Vars))
}
