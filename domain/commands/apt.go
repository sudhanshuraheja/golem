package commands

import (
	"fmt"
	"strings"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
)

type Apt struct {
	PGP              *string        `hcl:"pgp"`
	Repository       *APTRepository `hcl:"repository,block"`
	Update           *bool          `hcl:"update"`
	Purge            *[]string      `hcl:"purge"`
	Install          *[]string      `hcl:"install"`
	InstallNoUpgrade *[]string      `hcl:"install_no_upgrade"`
}

type APTRepository struct {
	URL     string `hcl:"url"`
	Sources string `hcl:"sources"`
}

func NewAPT() *Apt {
	return &Apt{}
}

func (a *Apt) Prepare() (Commands, artifacts.Artifacts) {
	cmds := Commands{}
	artfs := artifacts.Artifacts{}

	template := "#!/bin/bash\n"

	curl := ""
	pgp := ""
	if a.PGP != nil {
		curl = a.install("curl", false)
		template += a.ifPkgExists("curl", "", "", curl)

		pgp = a.pgp(*a.PGP)
	}

	common := ""
	repo := ""
	if a.Repository != nil && a.Repository.URL != "" && a.Repository.Sources != "" {
		common = a.install("software-properties-common", false)
		template += a.ifPkgExists("software-properties-common", "", "", common)

		repo = a.addRepository(a.Repository.URL, a.Repository.Sources)
	}

	pkgs := []string{}
	if a.Install != nil {
		for _, pkg := range *a.Install {
			if pkg != "" {
				install := a.install(pkg, false)
				template += a.ifPkgExists(pkg, pgp, repo, install)
				pgp = ""
				repo = ""
				pkgs = append(pkgs, pkg)
			}
		}
	}

	if a.InstallNoUpgrade != nil {
		for _, pkg := range *a.InstallNoUpgrade {
			if pkg != "" {
				install := a.install(pkg, true)
				template += a.ifPkgExists(pkg, pgp, repo, install)
				pgp = ""
				repo = ""
				pkgs = append(pkgs, pkg)
			}
		}
	}

	if a.Purge != nil {
		for _, pkg := range *a.Purge {
			if pkg != "" {
				purge := a.purge(pkg)
				template += purge
				pkgs = append(pkgs, pkg)
			}
		}
	}

	if a.Update != nil {
		template += a.update()
	}

	destination := fmt.Sprintf(
		"./temp/apt-%s-%s.sh",
		strings.Join(pkgs, "_"),
		utils.UUID().GetShort(),
	)
	artifact := artifacts.Artifact{
		Template: &artifacts.ArtifactTemplate{
			Data: &template,
		},
		Destination: destination,
	}
	artfs = append(artfs, &artifact)

	chmod := fmt.Sprintf("chmod 755 %s", destination)
	execute := destination
	remove := fmt.Sprintf("rm %s", destination)

	cmds.Append(NewCommand(chmod))
	cmds.Append(NewCommand(execute))
	cmds.Append(NewCommand(remove))

	return cmds, artfs
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
