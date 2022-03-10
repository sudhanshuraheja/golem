package commands

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestApt(t *testing.T) {
	trueValue := true
	dockerPGP := "https://download.docker.com/linux/ubuntu/gpg"
	dockerInstall := []string{"docker-ce", "docker-ce-cli", "containerd.io"}
	hashiPGP := "https://apt.releases.hashicorp.com/gpg"
	hashiInstall := []string{"nomad"}

	a1 := Apt{
		Update: &trueValue,
	}

	a2 := Apt{
		PGP: &dockerPGP,
		Repository: &APTRepository{
			URL:     "https://download.docker.com/linux/ubuntu",
			Sources: "stable",
		},
		Install: &dockerInstall,
	}

	a3 := Apt{
		PGP: &hashiPGP,
		Repository: &APTRepository{
			URL:     "https://apt.releases.hashicorp.com",
			Sources: "main",
		},
		Install: &hashiInstall,
	}

	commands, artifacts := a1.Prepare()
	utils.Test().Equals(t, 3, len(commands))
	utils.Test().Equals(t, 1, len(artifacts))
	utils.Test().Contains(t, *artifacts[0].Template.Data, "apt-get update")

	commands, artifacts = a2.Prepare()
	utils.Test().Equals(t, 3, len(commands))
	utils.Test().Equals(t, 1, len(artifacts))
	utils.Test().Contains(t, *artifacts[0].Template.Data, "docker-ce-cli")

	commands, artifacts = a3.Prepare()
	utils.Test().Equals(t, 3, len(commands))
	utils.Test().Equals(t, 1, len(artifacts))
	utils.Test().Contains(t, *artifacts[0].Template.Data, "nomad")

}
