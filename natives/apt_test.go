package natives

import (
	"testing"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

func TestApt(t *testing.T) {
	apt := NewAPT()

	a1 := config.Apt{}
	trueValue := true
	a1.Update = &trueValue

	a2 := config.Apt{}
	pgp := "https://download.docker.com/linux/ubuntu/gpg"
	a2.PGP = &pgp
	a2.Repository = &config.APTRepository{
		URL:     "https://download.docker.com/linux/ubuntu",
		Sources: "stable",
	}
	a2install := []string{"docker-ce", "docker-ce-cli", "containerd.io"}
	a2.Install = &a2install

	a3 := config.Apt{}
	pgp2 := "https://apt.releases.hashicorp.com/gpg"
	a3.PGP = &pgp2
	a3.Repository = &config.APTRepository{
		URL:     "https://apt.releases.hashicorp.com",
		Sources: "main",
	}
	a3install := []string{"nomad"}
	a3.Install = &a3install

	commands, artifacts := apt.Prepare([]config.Apt{
		a1,
		a2,
		a3,
	})
	utils.Test().Equals(t, 9, len(commands))
	utils.Test().Equals(t, 3, len(artifacts))
}
