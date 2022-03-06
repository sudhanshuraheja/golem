package natives

import (
	"fmt"
	"strings"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

type Apt struct {
}

func NewAPT() *Apt {
	return &Apt{}
}

func (a *Apt) ifPkgExists(pkg, pgp, repo, cmd string) string {
	installTemplate := `
pkg=%s
status="$(dpkg-query -W --showformat='${db:Status-Status}' "$pkg" 2>&1)"
if [ ! $? = 0 ] || [ ! "$status" = installed ]; then
	%s
	%s
	%s
fi
	`
	return fmt.Sprintf(installTemplate, pkg, pgp, repo, cmd)
}

func (a *Apt) update() string {
	return "sudo apt-get update --quiet --assume-yes\n"
}

func (a *Apt) purge(pkg string) string {
	format := "sudo apt-get purge %s --quiet --assume-yes\n"
	return fmt.Sprintf(format, pkg)
}

func (a *Apt) pgp(pgp string) string {
	format := "curl -fsSL %s | sudo apt-key add -\n"
	return fmt.Sprintf(format, pgp)
}

func (a *Apt) addRepository(url, sources string) string {
	format := "sudo apt-add-repository \"deb [arch=$(dpkg --print-architecture)] %s $(lsb_release -cs) %s\"\n"
	return fmt.Sprintf(format, url, sources)
}

func (a *Apt) install(pkg string, noUpgradeFlag bool) string {
	noUpgrade := ""
	if noUpgradeFlag {
		noUpgrade = "--no-upgrade"
	}
	format := "sudo apt-get install %s %s --quiet --assume-yes\n"
	return fmt.Sprintf(format, pkg, noUpgrade)
}

func (a *Apt) Prepare(capt []config.Apt) ([]config.Command, []config.Artifact) {
	commands := []config.Command{}
	artifacts := []config.Artifact{}

	for _, apt := range capt {
		template := "#!/bin/bash\n"

		curl := ""
		pgp := ""
		if apt.PGP != nil {
			curl = a.install("curl", false)
			template += a.ifPkgExists("curl", "", "", curl)

			pgp = a.pgp(*apt.PGP)
		}

		common := ""
		repo := ""
		if apt.Repository != nil {
			common = a.install("software-properties-common", false)
			template += a.ifPkgExists("software-properties-common", "", "", common)

			repo = a.addRepository(apt.Repository.URL, apt.Repository.Sources)
		}

		pkgs := []string{}
		if apt.Install != nil {
			for _, pkg := range *apt.Install {
				if pkg != "" {
					install := a.install(pkg, false)
					template += a.ifPkgExists(pkg, pgp, repo, install)
					pgp = ""
					repo = ""
					pkgs = append(pkgs, pkg)
				}
			}
		}

		if apt.InstallNoUpgrade != nil {
			for _, pkg := range *apt.InstallNoUpgrade {
				if pkg != "" {
					install := a.install(pkg, true)
					template += a.ifPkgExists(pkg, pgp, repo, install)
					pgp = ""
					repo = ""
					pkgs = append(pkgs, pkg)
				}
			}
		}

		if apt.Purge != nil {
			for _, pkg := range *apt.Purge {
				if pkg != "" {
					purge := a.purge(pkg)
					template += purge
					pkgs = append(pkgs, pkg)
				}
			}
		}

		if apt.Update != nil {
			template += a.update()
		}

		destination := fmt.Sprintf(
			"./temp/apt-%s-%s.sh",
			strings.Join(pkgs, "_"),
			utils.UUID().GetShort(),
		)
		artifact := config.Artifact{
			Template: &config.Template{
				Data: &template,
			},
			Destination: destination,
		}
		artifacts = append(artifacts, artifact)

		chmod := fmt.Sprintf("chmod 755 %s", destination)
		execute := destination
		remove := fmt.Sprintf("rm %s", destination)

		commands = append(commands, config.Command{Exec: &chmod})
		commands = append(commands, config.Command{Exec: &execute})
		commands = append(commands, config.Command{Exec: &remove})

	}
	return commands, artifacts
}
