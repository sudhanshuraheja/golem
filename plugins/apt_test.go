package plugins

import (
	"testing"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

func TestApt(t *testing.T) {
	apt := NewAPT()

	trueValue := true
	dockerPGP := "https://download.docker.com/linux/ubuntu/gpg"
	dockerInstall := []string{"docker-ce", "docker-ce-cli", "containerd.io"}
	hashiPGP := "https://apt.releases.hashicorp.com/gpg"
	hashiInstall := []string{"nomad"}

	a1 := config.Apt{
		Update: &trueValue,
	}

	a2 := config.Apt{
		PGP: &dockerPGP,
		Repository: &config.APTRepository{
			URL:     "https://download.docker.com/linux/ubuntu",
			Sources: "stable",
		},
		Install: &dockerInstall,
	}

	a3 := config.Apt{
		PGP: &hashiPGP,
		Repository: &config.APTRepository{
			URL:     "https://apt.releases.hashicorp.com",
			Sources: "main",
		},
		Install: &hashiInstall,
	}

	commands, artifacts := apt.Prepare([]config.Apt{a1, a2, a3})
	utils.Test().Equals(t, 9, len(commands))
	utils.Test().Equals(t, 3, len(artifacts))
}
