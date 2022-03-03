package natives

import (
	"fmt"
	"strings"

	"github.com/sudhanshuraheja/golem/config"
)

type Apt struct {
}

func NewAPT() *Apt {
	return &Apt{}
}

func (a *Apt) PGP(pgp string) (string, error) {
	if pgp == "" {
		return "", fmt.Errorf("pgp: pgp should not be empty")
	}
	format := "curl -fsSL %s | sudo apt-key add -"
	return fmt.Sprintf(format, pgp), nil
}

func (a *Apt) Repository(url, sources string) (string, error) {
	if url == "" || sources == "" {
		return "", fmt.Errorf("repository: url or sources should not be empty")
	}
	format := "sudo apt-add-repository \"deb [arch=$(dpkg --print-architecture)] %s $(lsb_release -cs) %s\""
	return fmt.Sprintf(format, url, sources), nil
}

func (a *Apt) Update() string {
	return "sudo apt-get update --quiet --assume-yes"
}

func (a *Apt) Purge(packages []string) (string, error) {
	format := "sudo apt-get purge %s"
	pkg, err := a.flattenStringArray(packages, " ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(format, pkg), nil
}

func (a *Apt) Install(packages []string) (string, error) {
	format := "sudo apt-get install %s --quiet --assume-yes"
	pkg, err := a.flattenStringArray(packages, " ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(format, pkg), nil
}

func (a *Apt) InstallNoUpgrade(packages []string) (string, error) {
	format := "sudo apt-get install %s --no-upgrade --quiet --assume-yes"
	pkg, err := a.flattenStringArray(packages, " ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(format, pkg), nil
}

func (a *Apt) flattenStringArray(packages []string, separator string) (string, error) {
	if len(packages) == 0 {
		return "", fmt.Errorf("packages cannot be empty")
	}
	merged := strings.Join(packages, separator)
	if merged == "" {
		return "", fmt.Errorf("merged packages is empty")
	}
	return merged, nil
}

func (a *Apt) ParseConfig(capt []config.Apt) ([]string, error) {
	commands := []string{}
	for _, apt := range capt {
		if apt.PGP != nil {
			// install curl
			cmd, err := a.Install([]string{"curl"})
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)

			// add pgp
			cmd, err = a.PGP(*apt.PGP)
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)
		}
		if apt.Repository != nil {
			// install software-properties-common
			cmd, err := a.Install([]string{"software-properties-common"})
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)

			// add repo
			cmd, err = a.Repository(apt.Repository.URL, apt.Repository.Sources)
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)
		}
		if apt.Update != nil {
			commands = append(commands, a.Update())
		}
		if apt.Purge != nil {
			cmd, err := a.Purge(*apt.Purge)
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)
		}
		if apt.Install != nil {
			cmd, err := a.Install(*apt.Install)
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)
		}
		if apt.InstallNoUpgrade != nil {
			cmd, err := a.InstallNoUpgrade(*apt.InstallNoUpgrade)
			if err != nil {
				return commands, err
			}
			commands = append(commands, cmd)
		}
	}
	return commands, nil
}
