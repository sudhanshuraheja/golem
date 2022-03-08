package recipes

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/betas-in/logger"
	"github.com/betas-in/pool"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/commands"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

var (
	tiny = 50
)

type Recipe struct {
	Name    string
	OfType  string
	Match   *servers.Match
	KV      map[string]string
	log     *logger.CLILogger
	tpl     *template.Template
	servers []servers.Server
	cmds    []commands.Command
	artfs   []artifacts.Artifact
}

func NewRecipe(log *logger.CLILogger, tpl *template.Template) *Recipe {
	return &Recipe{
		log: log,
		tpl: tpl,
		KV:  map[string]string{},
	}
}

func (r *Recipe) Execute(srvs []servers.Server, kv *kv.KV, maxParallelProcesses int) {
	r.processKV(kv)
	r.findServers(srvs)
	r.PrepareArtifacts(r.artfs, false)
	r.askPermission()

	r.run(maxParallelProcesses)
}

func (r *Recipe) AddCommand(cmd commands.Command) {
	r.cmds = append(r.cmds, cmd)
}

func (r *Recipe) AddArtifact(artf artifacts.Artifact) {
	r.artfs = append(r.artfs, artf)
}

func (r *Recipe) AddServers(srvs []servers.Server) {
	r.servers = srvs
}

func (r *Recipe) PrepareArtifacts(artfs []artifacts.Artifact, dryrun bool) {
	newArtfs := []artifacts.Artifact{}
	for _, a := range artfs {
		err := a.HandlePath(r.log, r.tpl)
		if err != nil {
			r.log.Error(r.Name).Msgf("coult not handle path: %v", err)
			continue
		}

		err = a.HandleData(r.tpl, dryrun)
		if err != nil {
			r.log.Error(r.Name).Msgf("could not handle data: %v", err)
			continue
		}

		err = a.HandleSource(r.tpl)
		if err != nil {
			r.log.Error(r.Name).Msgf("could not handle source: %v", err)
			continue
		}

		err = a.HandleDestination(r.tpl)
		if err != nil {
			r.log.Error(r.Name).Msgf("could not handle destination: %v", err)
			continue
		}

		newArtfs = append(newArtfs, a)
	}
	r.artfs = newArtfs
}

func (r *Recipe) processKV(k *kv.KV) {
	for key, value := range r.KV {
		existingValue, err := k.Get(key)
		if err != nil || existingValue == "" {
			err = k.Set(key, value)
			if err != nil {
				r.log.Error(r.Name).Msgf("could not set up kv: %s with value %s: %v", key, value, err)
				return
			}
			r.log.Info(r.Name).Msgf(
				"setup kv %s%s",
				logger.CyanBold("@golem.kv."),
				logger.CyanBold(key),
			)
		}
	}
}

func (r *Recipe) Display(query string) {
	name := logger.CyanBold(r.Name)
	if r.OfType == "local" {
		name = logger.Cyan(r.Name)
	}

	var attribute, operator, value string
	match := ""
	if r.Match != nil {
		attribute = r.Match.Attribute
		operator = r.Match.Operator
		value = r.Match.Value
		match = fmt.Sprintf("%s %s %s ", attribute, operator, value)
		match = logger.Yellow(match)
	}

	if len(query) > 0 {
		if !strings.Contains(name, query) {
			return
		}
	}

	r.log.Info(r.OfType).Msgf(
		"%s %s",
		name,
		match,
	)

	r.PrepareArtifacts(r.artfs, true)
	r.displayPrepared()
}

func (r *Recipe) findServers(all []servers.Server) {
	switch r.OfType {
	case "remote":
		if r.Match == nil {
			r.log.Error(r.Name).Msgf("recipe needs a 'match' block because of 'remote' execution")
			return
		}

		var err error
		if r.tpl != nil {
			r.Match.Value, err = r.tpl.Execute(r.Match.Value)
			if err != nil {
				r.log.Error(r.Name).Msgf("Error parsing template <%s>: %v", r.Match.Value, err)
				return
			}
		}

		r.servers, err = servers.NewMatch(r.Match.Attribute, r.Match.Operator, r.Match.Value).Find(all)
		if err != nil {
			r.log.Error(r.Name).Msgf("%v", err)
			return
		}
		serverNames := []string{}
		for _, s := range r.servers {
			serverNames = append(serverNames, s.Name)
		}

		if len(r.servers) == 0 {
			r.log.Highlight(r.Name).Msgf("no servers matched '%s %s %s'", r.Match.Attribute, r.Match.Operator, r.Match.Value)
			return
		}

		r.log.Info(r.Name).Msgf("found %d matching servers - %s", len(r.servers), strings.Join(serverNames, ", "))

	case "local":
	default:
		r.log.Error(r.Name).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipe) askPermission() {
	r.displayPrepared()
	answer := localutils.Question(r.log, r.Name, "Are you sure you want to continue [y/n]?")
	if utils.Array().Contains([]string{"y", "yes"}, answer, false) == -1 {
		r.log.Error(r.Name).Msgf("Quitting, because you said %s", answer)
		os.Exit(0)
	}
}

func (r *Recipe) run(maxParallelProcesses int) {
	r.downloadArtifacts()

	switch r.OfType {
	case "remote":
		if len(r.servers) == 0 {
			r.log.Error(r.Name).Msgf("no matching servers found")
			return
		}

		r.startSSHPool(int64(maxParallelProcesses))
	case "local":
		c := Cmd{log: r.log}
		c.Upload(r.artfs)
		c.Run(r.cmds)
	default:
		r.log.Error(r.Name).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipe) downloadArtifacts() {
	for i, a := range r.artfs {
		source, err := a.Download(r.log)
		if err != nil {
			r.log.Error(r.Name).Msgf("%v", err)
			os.Exit(1)
		}
		r.artfs[i].Source = &source
	}
}

func (r *Recipe) startSSHPool(maxProcs int64) {
	log := logger.NewLogger(2, true)

	wp := pool.NewPool("ssh", log)
	wp.AddWorkerGroup(NewSSHWorkerGroup("ssh", r.log, 5*time.Second))

	processed := wp.Start(maxProcs)

	startTime := time.Now()
	for _, s := range r.servers {
		wp.Queue(SSHJob{Server: s, Recipe: r})
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
			if count == len(r.servers) {
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

	r.log.Announce(r.Name).Msgf("completed %s", localutils.TimeInSecs(startTime))
}

func (r *Recipe) displayPrepared() {
	for _, ar := range r.artfs {
		sourcePath := ""
		switch {
		case ar.Template.Data != nil:
			sourcePath = *ar.Template.Data
		case ar.Template.Path != nil:
			sourcePath = *ar.Template.Path
		case ar.Source != nil:
			sourcePath = *ar.Source
		}

		r.log.Info("").Msgf(
			"%s %s %s %s",
			logger.Cyan("uploading"),
			localutils.TinyString(sourcePath, tiny),
			logger.Cyan("to"),
			localutils.TinyString(ar.Destination, tiny),
		)
	}

	for _, command := range r.cmds {
		exec, err := r.tpl.Execute(command.Exec)
		if err != nil {
			r.log.Error(r.Name).Msgf("could not parse template %s: %v", command.Exec, err)
		}
		r.log.Info(r.Name).Msgf("$ %s", localutils.TinyString(exec, tiny*2))
	}
}
