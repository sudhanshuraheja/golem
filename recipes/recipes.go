package recipes

import (
	"fmt"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Recipes struct {
	conf *config.Config
	log  *logger.CLILogger
	tpl  *Template
}

func NewRecipes(conf *config.Config, log *logger.CLILogger) *Recipes {
	r := Recipes{
		conf: conf,
		log:  log,
	}
	r.tpl = &Template{}
	if conf.Vars != nil {
		r.tpl.Vars = *conf.Vars
	}
	r.tpl.Servers = append(r.tpl.Servers, conf.Servers...)
	return &r
}

func (r *Recipes) List(query string) {
	r.log.Announce("").Msgf("list of all available recipes")

	// Add system defined
	r.log.Info("system").Msgf("%s", logger.Cyan("list"))
	r.log.Info("system").Msgf("%s", logger.Cyan("servers"))

	for _, re := range r.conf.Recipes {
		recipe := Recipe{log: r.log, base: &re}

		name := logger.CyanBold(re.Name)
		if re.Type == "local" {
			name = logger.Cyan(re.Name)
		}

		var attribute, operator, value string
		match := ""
		if re.Match != nil {
			attribute = re.Match.Attribute
			operator = re.Match.Operator
			value = re.Match.Value
			match = fmt.Sprintf("%s %s %s ", attribute, operator, value)
			match = logger.Yellow(match)
		}

		if len(query) > 0 {
			if !strings.Contains(name, query) {
				continue
			}
		}

		r.log.Info(re.Type).Msgf(
			"%s %s",
			name,
			match,
		)

		recipe.PrepareArtifacts(r.tpl, true)

		for _, ar := range recipe.preparedArtifacts {

			if ar.Template != nil {
				r.log.Info("").Msgf("%s %s %s %s", logger.Cyan("uploading"), *ar.Template, logger.Cyan("to"), ar.Destination)
			} else {
				r.log.Info("").Msgf("%s %s %s %s", logger.Cyan("uploading"), ar.Source, logger.Cyan("to"), ar.Destination)
			}
		}

		recipe.PrepareCommands(r.tpl)

		for _, command := range recipe.preparedCommands {
			r.log.Info("").Msgf("$ %s", command)
		}
	}
}

func (r *Recipes) Servers(query string) {
	r.log.Announce("").Msgf("list of all connected servers")

	for _, s := range r.conf.Servers {

		if len(query) > 0 {
			if !strings.Contains(s.Name, query) {
				continue
			}
		}

		primaryName := s.Name
		if primaryName == "" {
			primaryName = localutils.StringPtrValue(s.PublicIP, "")
		}

		username := ""
		if s.User != "" {
			username = fmt.Sprintf("%s %s ", logger.Cyan("user"), s.User)
		}

		port := ""
		if s.Port != 0 {
			port = fmt.Sprintf("%s %d ", logger.Cyan("port"), s.Port)
		}

		publicIP := localutils.StringPtrValue(s.PublicIP, "")
		if publicIP != "" {
			publicIP = fmt.Sprintf("%s %s ", logger.Cyan("publicIP"), publicIP)
		}

		privateIP := localutils.StringPtrValue(s.PrivateIP, "")
		if privateIP != "" {
			privateIP = fmt.Sprintf("%s %s ", logger.Cyan("privateIP"), privateIP)
		}

		r.log.Info(primaryName).Msgf(
			"%s%s%s%s",
			username,
			port,
			publicIP,
			privateIP,
		)

		hostnames := strings.Join(s.HostName, ", ")
		if hostnames != "" {
			r.log.Info("").Msgf("%s %s", logger.Cyan("hosts"), hostnames)
		}

		tags := strings.Join(s.Tags, ", ")
		if tags != "" {
			r.log.Info("").Msgf("%s %s", logger.Cyan("tags"), tags)
		}
	}
}

func (r *Recipes) Run(name string) {
	recipe := Recipe{log: r.log}

	for i, re := range r.conf.Recipes {
		if re.Name == name {
			recipe.base = &r.conf.Recipes[i]
		}
	}

	if recipe.base.Name == "" {
		r.log.Error(name).Msgf("the recipe <%s> was not found in '~/.golem/' or '.'", name)
		return
	}

	recipe.FindServers(r.conf.Servers)

	recipe.PrepareArtifacts(r.tpl, false)
	recipe.PrepareCommands(r.tpl)

	recipe.AskPermission()

	recipe.DownloadArtifacts()
	recipe.ExecuteCommands(r.conf.MaxParallelProcesses)
}
