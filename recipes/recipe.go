package recipes

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/betas-in/getter"
	"github.com/betas-in/logger"
	"github.com/betas-in/pool"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Recipes struct {
	conf *config.Config
	log  *logger.CLILogger
}

func NewRecipes(conf *config.Config, log *logger.CLILogger) *Recipes {
	return &Recipes{conf: conf, log: log}
}

func (r *Recipes) List() {
	r.log.Announce("").Msgf("list of all available recipes")

	// Add system defined
	r.log.Info("system").Msgf("%s", logger.Cyan("list"))
	r.log.Info("system").Msgf("%s", logger.Cyan("servers"))

	for _, re := range r.conf.Recipes {
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

		r.log.Info(re.Type).Msgf(
			"%s %s",
			name,
			match,
		)

		for _, ar := range re.Artifacts {
			r.log.Info("").Msgf("%s %s %s", ar.Source, logger.Cyan("to"), ar.Destination)
		}

		for _, cm := range re.CustomCommands {
			r.log.Info("").Msgf("%s %s", logger.Cyan("$"), strings.TrimSuffix(cm.Exec, "\n"))
		}

		if re.Commands != nil {
			for _, cm := range *re.Commands {
				r.log.Info("").Msgf("%s %s", logger.Cyan("$"), cm)
			}
		}
	}
}

func (r *Recipes) Servers() {
	r.log.Announce("").Msgf("list of all connected servers")

	for _, s := range r.conf.Servers {

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

		hostnames := localutils.StringPtrValue(s.HostName, "")
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
	var recipe config.Recipe

	for i, re := range r.conf.Recipes {
		if re.Name == name {
			recipe = r.conf.Recipes[i]
		}
	}

	if recipe.Name == "" {
		r.log.Error(name).Msgf("the recipe <%s> was not found in '~/.golem/' or '.'", name)
		return
	}

	servers := r.askPermission(&recipe)

	err := r.downloadRemoteArtifacts(&recipe)
	if err != nil {
		r.log.Error(name).Msgf("%v", err)
	}

	switch recipe.Type {
	case "remote":
		r.RemotePool(servers, recipe, *r.conf.MaxParallelProcesses)
	case "local":
		c := Cmd{log: r.log}
		if len(recipe.CustomCommands) > 0 {
			commands := []string{}
			for _, cmd := range recipe.CustomCommands {
				commands = append(commands, strings.TrimSuffix(cmd.Exec, "\n"))
			}
			c.Run(commands)
		}
		if recipe.Commands != nil {
			c.Run(*recipe.Commands)
		}
	default:
		r.log.Error(name).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipes) askPermission(recipe *config.Recipe) []config.Server {
	var servers []config.Server
	switch recipe.Type {
	case "remote":
		if recipe.Match == nil {
			r.log.Error(recipe.Name).Msgf("recipe needs a 'match' block because of 'remote' execution")
			return servers
		}

		var err error
		servers, err = NewMatch(*recipe.Match).Find(r.conf)
		if err != nil {
			r.log.Error(recipe.Name).Msgf("%v", err)
			return servers
		}
		serverNames := []string{}
		for _, s := range servers {
			serverNames = append(serverNames, s.Name)
		}

		if len(servers) == 0 {
			r.log.Highlight(recipe.Name).Msgf("no servers matched '%s %s %s'", recipe.Match.Attribute, recipe.Match.Operator, recipe.Match.Value)
			return servers
		}

		r.log.Info(recipe.Name).Msgf("found %d matching servers - %s", len(servers), strings.Join(serverNames, ", "))

	case "local":
	default:
		r.log.Error(recipe.Name).Msgf("recipe only supports ['remote', 'local'] types")
	}

	for _, a := range recipe.Artifacts {
		r.log.Info(recipe.Name).Msgf("%s %s %s %s", logger.Cyan("uploading"), a.Source, logger.Cyan("to"), a.Destination)
	}

	for _, cmd := range recipe.CustomCommands {
		r.log.Info(recipe.Name).Msgf("$ %s", strings.TrimSuffix(cmd.Exec, "\n"))
	}

	if recipe.Commands != nil {
		for _, c := range *recipe.Commands {
			parsed, err := ParseTemplate(c, r.conf)
			if err != nil {
				r.log.Error(recipe.Name).Msgf("Error parsing template <%s>: %v", c, err)
			}
			r.log.Info(recipe.Name).Msgf("$ %s", parsed)
		}
	}

	answer := localutils.Question(r.log, recipe.Name, "Are you sure you want to continue [y/n]?")
	if utils.Array().Contains([]string{"y", "yes"}, answer, false) == -1 {
		r.log.Error(recipe.Name).Msgf("Quitting, because you said %s", answer)
		os.Exit(0)
	}

	return servers
}

func (r *Recipes) downloadRemoteArtifacts(recipe *config.Recipe) error {
	for i, a := range recipe.Artifacts {
		if strings.HasPrefix(a.Source, "http://") || strings.HasPrefix(a.Source, "https://") {
			r.log.Info(recipe.Name).Msgf("%s %s", logger.Cyan("downloading"), a.Source)
			log := logger.NewLogger(3, true)
			g := getter.NewGetter(log)

			startTime := time.Now()
			response := g.FetchResponse(getter.Request{
				Path:       a.Source,
				SaveToDisk: true,
			})

			if response.Error != nil {
				return fmt.Errorf("could not download %s: %v", a.Source, response.Error)
			}
			if response.Code != 200 {
				return fmt.Errorf("received error code for %s: %d", a.Source, response.Code)
			}

			r.log.Highlight(recipe.Name).Msgf("downloaded %s to %s %s", a.Source, response.DataPath, localutils.TimeInSecs(startTime))
			recipe.Artifacts[i].Source = response.DataPath
		}
	}
	return nil
}

func (r *Recipes) RemotePool(servers []config.Server, recipe config.Recipe, maxProcs int) {
	log := logger.NewLogger(2, true)
	wp := pool.NewPool("ssh", log)
	wp.AddWorkerGroup(NewSSHWorkerGroup("ssh", r.log, 5*time.Second))
	processed := wp.Start(int64(maxProcs))

	startTime := time.Now()
	for _, s := range servers {
		wp.Queue(SSHJob{Server: s, Recipe: recipe})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	count := 0

	loop := true
	for loop {
		select {
		case <-processed:
			count++
			if count == len(servers) {
				wp.Stop(ctx)
				loop = false
				break
			}
		case <-quit:
			wp.Stop(ctx)
			loop = false
			break
		}
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	ticks := 0
	for ; true; <-ticker.C {
		ticks++
		if ticks >= 20 {
			break
		}
		if wp.GetWorkerCount() == 0 {
			break
		}
	}

	r.log.Announce(recipe.Name).Msgf("completed %s", localutils.TimeInSecs(startTime))
}
