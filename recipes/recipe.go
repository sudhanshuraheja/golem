package recipes

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/betas-in/logger"
	"github.com/betas-in/pool"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/natives"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Recipe struct {
	base              *config.Recipe
	log               *logger.CLILogger
	servers           []config.Server
	preparedCommands  []string
	preparedArtifacts []config.Artifact
}

func (r *Recipe) FindServers(servers []config.Server) {
	switch r.base.Type {
	case "remote":
		if r.base.Match == nil {
			r.log.Error(r.base.Name).Msgf("recipe needs a 'match' block because of 'remote' execution")
			return
		}

		var err error
		r.servers, err = NewMatch(*r.base.Match).Find(servers)
		if err != nil {
			r.log.Error(r.base.Name).Msgf("%v", err)
			return
		}
		serverNames := []string{}
		for _, s := range r.servers {
			serverNames = append(serverNames, s.Name)
		}

		if len(r.servers) == 0 {
			r.log.Highlight(r.base.Name).Msgf("no servers matched '%s %s %s'", r.base.Match.Attribute, r.base.Match.Operator, r.base.Match.Value)
			return
		}

		r.log.Info(r.base.Name).Msgf("found %d matching servers - %s", len(r.servers), strings.Join(serverNames, ", "))

	case "local":
	default:
		r.log.Error(r.base.Name).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipe) PrepareArtifacts(tpl *Template, dryrun bool) {
	for _, a := range r.base.Artifacts {
		artifact := config.Artifact{}

		if a.Template != nil {

			if a.Template.Path != nil {
				if strings.HasPrefix(*a.Template.Path, "http://") || strings.HasPrefix(*a.Template.Path, "https://") {
					// Url based template
					path, err := Download(r.log, r.base.Name, *a.Template.Path)
					if err != nil {
						r.log.Error(r.base.Name).Msgf("%v", err)
						os.Exit(1)
					}
					a.Template.Path = &path
				} // else File base template

				bytes, err := os.ReadFile(*a.Template.Path)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
					os.Exit(1)
				}
				bytesString := string(bytes)
				a.Template.Data = &bytesString
			}

			if a.Template.Data != nil {
				parsedTemplate, err := ParseTemplate(*a.Template.Data, tpl)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("Error parsing template <%s>: %v", a.Template.Data, err)
					continue
				}
				artifact.Template = &config.Template{
					Data: &parsedTemplate,
				}

				if !dryrun {
					fileName, err := localutils.FileCopy(parsedTemplate)
					if err != nil {
						r.log.Error(r.base.Name).Msgf("could not save file: %v", err)
						os.Exit(1)
					}
					artifact.Source = &fileName
				}
			}
		}

		if a.Source != nil {
			parsedSource, err := ParseTemplate(*a.Source, tpl)
			if err != nil {
				r.log.Error(r.base.Name).Msgf("Error parsing template <%s>: %v", *a.Source, err)
				continue
			}
			artifact.Source = &parsedSource
		}

		parsedDestination, err := ParseTemplate(a.Destination, tpl)
		if err != nil {
			r.log.Error(r.base.Name).Msgf("Error parsing template <%s>: %v", a.Destination, err)
			continue
		}
		artifact.Destination = parsedDestination

		r.preparedArtifacts = append(r.preparedArtifacts, artifact)
	}
}

func (r *Recipe) PrepareCommands(tpl *Template) {
	for _, cmd := range r.base.CustomCommands {
		if cmd.Exec != nil {
			parsedCmd, err := ParseTemplate(*cmd.Exec, tpl)
			if err != nil {
				r.log.Error(r.base.Name).Msgf("Error parsing template <%s>: %v", *cmd.Exec, err)
				continue
			}
			parsedCmd = strings.TrimSuffix(parsedCmd, "\n")
			r.AddPreparedCommand(parsedCmd)
		}

		apt := natives.NewAPT()
		commands, err := apt.ParseConfig(cmd.Apt)
		if err != nil {
			r.log.Error(r.base.Name).Msgf("Error parsing apt: %v", err)
			continue
		}
		for _, cmd := range commands {
			r.AddPreparedCommand(cmd)
		}
	}

	if r.base.Commands != nil {
		for _, c := range *r.base.Commands {
			parsedCmd, err := ParseTemplate(c, tpl)
			if err != nil {
				r.log.Error(r.base.Name).Msgf("Error parsing template <%s>: %v", c, err)
				continue
			}
			parsedCmd = strings.TrimSuffix(parsedCmd, "\n")
			r.preparedCommands = append(r.preparedCommands, parsedCmd)
		}
	}
}

func (r *Recipe) AddPreparedCommand(cmd string) {
	r.preparedCommands = append(r.preparedCommands, cmd)
}

func (r *Recipe) AskPermission() {
	for _, a := range r.preparedArtifacts {
		if a.Template != nil {
			if a.Template.Data != nil {
				r.log.Info(r.base.Name).Msgf("%s %s %s %s", logger.Cyan("uploading"), *a.Template.Data, logger.Cyan("to"), a.Destination)
			}
			if a.Template.Path != nil {
				r.log.Info(r.base.Name).Msgf("%s %s %s %s", logger.Cyan("uploading"), *a.Template.Path, logger.Cyan("to"), a.Destination)
			}
		} else {
			r.log.Info(r.base.Name).Msgf("%s %s %s %s", logger.Cyan("uploading"), *a.Source, logger.Cyan("to"), a.Destination)
		}
	}

	for _, command := range r.preparedCommands {
		r.log.Info(r.base.Name).Msgf("$ %s", command)
	}

	answer := localutils.Question(r.log, r.base.Name, "Are you sure you want to continue [y/n]?")
	if utils.Array().Contains([]string{"y", "yes"}, answer, false) == -1 {
		r.log.Error(r.base.Name).Msgf("Quitting, because you said %s", answer)
		os.Exit(0)
	}
}

func (r *Recipe) Execute(maxParallelProcesses *int) {
	r.DownloadArtifacts()

	switch r.base.Type {
	case "remote":
		if len(r.servers) == 0 {
			r.log.Error(r.base.Name).Msgf("no matching servers found")
			return
		}

		maxProcs := 4
		if maxParallelProcesses != nil {
			maxProcs = *maxParallelProcesses
		}

		r.StartSSHPool(int64(maxProcs))
	case "local":
		c := Cmd{log: r.log}
		c.Upload(r.preparedArtifacts)
		c.Run(r.preparedCommands)
	default:
		r.log.Error(r.base.Name).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipe) DownloadArtifacts() {
	for i, a := range r.preparedArtifacts {
		if a.Source != nil && strings.HasPrefix(*a.Source, "http://") || strings.HasPrefix(*a.Source, "https://") {

			filePath, err := Download(r.log, r.base.Name, *a.Source)
			if err != nil {
				r.log.Error(r.base.Name).Msgf("%v", err)
				os.Exit(1)
			}

			r.base.Artifacts[i].Source = &filePath
		}
	}
}

func (r *Recipe) StartSSHPool(maxProcs int64) {
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

	r.log.Announce(r.base.Name).Msgf("completed %s", localutils.TimeInSecs(startTime))
}
