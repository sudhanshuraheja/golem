package natives

import (
	"fmt"
	"strings"
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
