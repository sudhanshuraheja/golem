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

var (
	header = "\n➡️  "
)

type Recipes struct {
	conf *config.Config
}

func NewRecipes(conf *config.Config) *Recipes {
	return &Recipes{conf: conf}
}

func (r *Recipes) List() {
	logger.Announcef("%sRecipes", header)
	tb := logger.NewCLITable("Name", "Match", "Artifacts", "Commands")
	for _, r := range r.conf.Recipes {
		var attribute, operator, value string
		if r.Match != nil {
			attribute = r.Match.Attribute
			operator = r.Match.Operator
			value = r.Match.Value
		}
		tb.Row(
			r.Name,
			fmt.Sprintf("%s %s %s", attribute, operator, value),
			len(r.Artifacts),
			len(r.Commands),
		)
	}
	// Add system defined
	tb.Row("servers", "local only", 0, 0)
	tb.Display()
}

func (r *Recipes) Servers() {
	logger.Announcef("%sServers", header)
	t := logger.NewCLITable("Name", "Public IP", "Private IP", "User", "Port", "Tags", "Hostname")
	for _, s := range r.conf.Servers {
		hostnames := localutils.StringPtrValue(s.HostName, "")
		if len(hostnames) > 60 {
			hostnames = hostnames[:60]
		}
		t.Row(
			s.Name,
			localutils.StringPtrValue(s.PublicIP, ""),
			localutils.StringPtrValue(s.PrivateIP, ""),
			s.User,
			s.Port,
			strings.Join(s.Tags, ", "),
			hostnames,
		)
	}
	t.Display()
}

func (r *Recipes) Run(name string) {
	var recipe config.Recipe

	for i, re := range r.conf.Recipes {
		if re.Name == name {
			recipe = r.conf.Recipes[i]
		}
	}

	if recipe.Name == "" {
		logger.Errorf("kitchen | the recipe <%s> was not found in '~/.golem/' or '.'", recipe.Name)
		return
	}

	servers := r.askPermission(&recipe)

	err := r.downloadRemoteArtifacts(&recipe)
	if err != nil {
		logger.Errorf("kitchen | %v", err)
	}

	switch recipe.Type {
	case "remote-exec":
		r.RemotePool(servers, recipe, *r.conf.MaxParallelProcesses)
	case "local-exec":
		c := Cmd{}
		c.Run(recipe.Commands)
	default:
		logger.Errorf("recipe only supports ['remote-exec', 'local-exec'] types")
	}
}

func (r *Recipes) askPermission(recipe *config.Recipe) []config.Server {
	var servers []config.Server
	switch recipe.Type {
	case "remote-exec":
		if recipe.Match == nil {
			logger.Errorf("kitchen | recipe <%s> need a 'match' block because of 'remote-exec'", recipe.Name)
			return servers
		}

		servers = NewMatch(*recipe.Match).Find(r.conf)
		serverNames := []string{}
		for _, s := range servers {
			serverNames = append(serverNames, s.Name)
		}

		if len(servers) == 0 {
			logger.MinorSuccessf("%s | no servers matched '%s %s %s'", recipe.Name, recipe.Match.Attribute, recipe.Match.Operator, recipe.Match.Value)
			return servers
		}

		logger.Announcef("%s | found %d matching servers - %s", recipe.Name, len(servers), strings.Join(serverNames, ", "))

	case "local-exec":
	default:
		logger.Errorf("recipe only supports ['remote-exec', 'local-exec'] types")
	}

	for _, a := range recipe.Artifacts {
		logger.Infof("→ %s → %s", a.Source, a.Destination)
	}

	for _, c := range recipe.Commands {
		parsed, err := ParseTemplate(c, r.conf)
		if err != nil {
			logger.Errorf("Error parsing template <%s>: %v", c, err)
		}
		logger.Infof("→ $ %s", parsed)
	}

	answer := logger.Questionf("Are you sure you want to continue [y/n]?")
	if utils.Array().Contains([]string{"y", "yes"}, answer, false) == -1 {
		logger.Errorf("Quitting, because you said %s", answer)
		os.Exit(0)
	}

	return servers
}

func (r *Recipes) downloadRemoteArtifacts(recipe *config.Recipe) error {
	for i, a := range recipe.Artifacts {
		if strings.HasPrefix(a.Source, "http://") || strings.HasPrefix(a.Source, "https://") {
			logger.Announcef("kitchen | downloading %s", a.Source)
			g := getter.NewGetter(nil)

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

			logger.MinorSuccessf("kitchen | downloaded %s to %s in %s", a.Source, response.DataPath, time.Since(startTime))
			recipe.Artifacts[i].Source = response.DataPath
		}
	}
	return nil
}

func (r *Recipes) RemotePool(servers []config.Server, recipe config.Recipe, maxProcs int) {
	log := logger.NewLogger(2, true)
	wp := pool.NewPool("ssh", log)
	wp.AddWorkerGroup(NewSSHWorkerGroup("ssh", 10*time.Millisecond))
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

	logger.MinorSuccessf("%s | completed in %s", recipe.Name, time.Since(startTime))
}
