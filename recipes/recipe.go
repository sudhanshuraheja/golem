package recipes

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/betas-in/getter"
	"github.com/betas-in/logger"
	"github.com/betas-in/pool"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/natives"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Recipe struct {
	base             *config.Recipe
	log              *logger.CLILogger
	servers          []config.Server
	preparedCommands []string
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

func (r *Recipe) PrepareCommands(tpl *Template) {
	for _, cmd := range r.base.CustomCommands {
		if cmd.Exec != nil {
			parsedCmd, err := ParseTemplate(*cmd.Exec, tpl)
			if err != nil {
				r.log.Error(r.base.Name).Msgf("Error parsing template <%s>: %v", cmd.Exec, err)
				continue
			}
			parsedCmd = strings.TrimSuffix(parsedCmd, "\n")
			r.AddPreparedCommand(parsedCmd)
		}

		aptCmd := natives.NewAPT()
		for _, apt := range cmd.Apt {
			if apt.PGP != nil {
				// install curl
				cmd, err := aptCmd.Install([]string{"curl"})
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)

				// add pgp
				cmd, err = aptCmd.PGP(*apt.PGP)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)
			}
			if apt.Repository != nil {
				// install software-properties-common
				cmd, err := aptCmd.Install([]string{"software-properties-common"})
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)

				// add repo
				cmd, err = aptCmd.Repository(apt.Repository.URL, apt.Repository.Sources)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)
			}
			if apt.Update != nil {
				r.AddPreparedCommand(aptCmd.Update())
			}
			if apt.Purge != nil {
				cmd, err := aptCmd.Purge(*apt.Purge)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)
			}
			if apt.Install != nil {
				cmd, err := aptCmd.Purge(*apt.Install)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)
			}
			if apt.InstallNoUpgrade != nil {
				cmd, err := aptCmd.Purge(*apt.InstallNoUpgrade)
				if err != nil {
					r.log.Error(r.base.Name).Msgf("%v", err)
				}
				r.AddPreparedCommand(cmd)
			}
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
	for _, a := range r.base.Artifacts {
		r.log.Info(r.base.Name).Msgf("%s %s %s %s", logger.Cyan("uploading"), a.Source, logger.Cyan("to"), a.Destination)
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

func (r *Recipe) DownloadArtifacts() {
	for i, a := range r.base.Artifacts {
		if strings.HasPrefix(a.Source, "http://") || strings.HasPrefix(a.Source, "https://") {
			r.log.Info(r.base.Name).Msgf("%s %s", logger.Cyan("downloading"), a.Source)
			log := logger.NewLogger(3, true)
			g := getter.NewGetter(log)

			startTime := time.Now()
			response := g.FetchResponse(getter.Request{
				Path:       a.Source,
				SaveToDisk: true,
			})

			if response.Error != nil {
				r.log.Error(r.base.Name).Msgf("could not download %s: %v", a.Source, response.Error)
				os.Exit(1)
			}
			if response.Code != 200 {
				r.log.Error(r.base.Name).Msgf("received error code for %s: %d", a.Source, response.Code)
				os.Exit(1)
			}

			r.log.Highlight(r.base.Name).Msgf("downloaded %s to %s %s", a.Source, response.DataPath, localutils.TimeInSecs(startTime))
			r.base.Artifacts[i].Source = response.DataPath
		}
	}
}

func (r *Recipe) ExecuteCommands(maxParallelProcesses *int) {
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
		c.Run(r.preparedCommands)
	default:
		r.log.Error(r.base.Name).Msgf("recipe only supports ['remote', 'local'] types")
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
