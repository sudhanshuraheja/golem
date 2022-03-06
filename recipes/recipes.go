package recipes

import (
	"fmt"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/kv"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
	"github.com/sudhanshuraheja/golem/template"
)

var (
	tiny = 50
)

type Recipes struct {
	conf *config.Config
	log  *logger.CLILogger
	tpl  *template.Template
	kv   *kv.KV
}

func NewRecipes(conf *config.Config, log *logger.CLILogger) *Recipes {
	r := Recipes{
		conf: conf,
		log:  log,
	}
	r.kv = kv.NewKV(log)
	r.tpl = template.NewTemplate(conf, r.kv)
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
				if ar.Template.Data != nil {
					r.log.Info("").Msgf(
						"%s %s %s %s",
						logger.Cyan("uploading"),
						localutils.TinyString(*ar.Template.Data, tiny),
						logger.Cyan("to"),
						localutils.TinyString(ar.Destination, tiny),
					)
				}
				if ar.Template.Path != nil {
					r.log.Info("").Msgf(
						"%s %s %s %s",
						logger.Cyan("uploading"),
						localutils.TinyString(*ar.Template.Path, tiny),
						logger.Cyan("to"),
						localutils.TinyString(ar.Destination, tiny),
					)
				}
			} else {
				r.log.Info("").Msgf(
					"%s %s %s %s",
					logger.Cyan("uploading"),
					localutils.TinyString(*ar.Source, tiny),
					logger.Cyan("to"),
					localutils.TinyString(ar.Destination, tiny),
				)
			}
		}

		recipe.PrepareCommands(r.tpl)

		for _, command := range recipe.preparedCommands {
			if command.Exec == nil {
				continue
			}
			r.log.Info("").Msgf("$ %s", *command.Exec)
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

func (r *Recipes) KV(path, action string) {
	if path == "" {
		path = "list"
	}
	if path == "list" {
		store, err := r.kv.GetAll()
		if err != nil {
			r.log.Error("kv").Msgf("could not read from the database: %v", err)
			_ = r.kv.Close()
			return
		}

		kvLog := func(log *logger.CLILogger, key, value string) {
			if strings.Contains(key, "secret") || strings.Contains(key, "password") {
				value = value[:2] + "************"
			}
			log.Info("kv").Msgf("%s: %s", logger.Cyan(key), logger.GreenBold(value))
		}

		for key, value := range store {
			switch action {
			case "":
				kvLog(r.log, key, value)
			default:
				if strings.Contains(key, action) {
					kvLog(r.log, key, value)
				}
			}
		}
		_ = r.kv.Close()
		return
	}

	switch action {
	case "set":
		userValue := localutils.Question(r.log, "enter", "please enter a value")
		userValue = strings.TrimSuffix(userValue, "\n")
		err := r.kv.Set(path, userValue)
		if err != nil {
			r.log.Error("kv").Msgf("could not set the value: %v", err)
		}
	case "rand32":
		err := r.kv.Set(path, action)
		if err != nil {
			r.log.Error("kv").Msgf("could not set the value: %v", err)
		}

		value, err := r.kv.Get(path)
		if err != nil {
			r.log.Error("kv").Msgf("could not get value: %v", err)
		}

		if err == nil {
			r.log.Info("kv").Msgf("%s: %s", logger.Cyan(path), logger.GreenBold(value))
		}
	case "delete":
		err := r.kv.Delete(path)
		if err != nil {
			r.log.Error("kv").Msgf("could not delete the value: %v", err)
		}
	default:
		value, err := r.kv.Get(path)
		if err != nil {
			r.log.Error("kv").Msgf("could not get value: %v", err)
		}

		if err == nil {
			r.log.Info("kv").Msgf("%s: %s", logger.Cyan(path), logger.GreenBold(value))
		}
	}

	_ = r.kv.Close()
}

func (r *Recipes) Run(name string) {
	recipe := Recipe{log: r.log}

	for i, re := range r.conf.Recipes {
		if re.Name == name {
			recipe.base = &r.conf.Recipes[i]
		}
	}

	if recipe.base == nil || recipe.base.Name == "" {
		r.log.Error(name).Msgf("the recipe %s was not found in '~/.golem/' or '.'", logger.Cyan(name))
		return
	}

	recipe.SetupKV(r.kv)

	recipe.FindServers(r.conf.Servers, r.tpl)

	recipe.PrepareArtifacts(r.tpl, false)
	recipe.PrepareCommands(r.tpl)

	recipe.AskPermission()

	recipe.Execute(r.conf.MaxParallelProcesses)
}
