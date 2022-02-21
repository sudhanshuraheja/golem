package recipes

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/pool"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

type SSHJob struct {
	Recipe config.Recipe
	Server config.Server
}

func List(c *config.Config) {
	tb := log.NewTable("Name", "Match", "Artifacts", "Commands")
	for _, r := range c.Recipe {
		tb.Row(
			r.Name,
			fmt.Sprintf("%s %s %s", r.Match.Attribute, r.Match.Operator, r.Match.Value),
			len(r.Artifacts),
			len(r.Commands),
		)
	}
	// Add system defined
	tb.Row("servers", "local only", 0, 0)
	tb.Display()
}

func Exists(c *config.Config, name string) bool {
	for _, r := range c.Recipe {
		if r.Name == name {
			return true
		}
	}
	return false
}

func Init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("init | could not find user's home directory: %v", err)
		return
	}
	confDir := fmt.Sprintf("%s/.golem", dirname)

	err = os.MkdirAll(confDir, os.FileMode(0755))
	if err != nil {
		log.Errorf("init | could not create conf dir %s: %v", confDir, err)
		return
	}

	confFile := fmt.Sprintf("%s/golem.hcl", confDir)
	_, err = os.Stat(confFile)
	if os.IsNotExist(err) {
		file, err := os.Create(confFile)
		if err != nil {
			log.Errorf("init | error creating conf file %s: %v", confFile, err)
			return
		}
		defer file.Close()
		log.MinorSuccessf("init | conf file created at %s", confFile)
	} else if err != nil {
		log.Errorf("init | error checking conf file %s: %v", confFile, err)
	}
	// log.MinorSuccessf("init | conf file already exists at %s", confFile)
}

func Run(c *config.Config, name string) {
	var recipe config.Recipe
	for i, r := range c.Recipe {
		if r.Name == name {
			recipe = c.Recipe[i]
		}
	}

	servers := findMatchingServers(c, recipe.Match)
	serverNames := []string{}
	for _, s := range servers {
		serverNames = append(serverNames, s.Name)
	}

	if len(servers) == 0 {
		log.MinorSuccessf("%s | no servers matched '%s %s %s'", recipe.Name, recipe.Match.Attribute, recipe.Match.Operator, recipe.Match.Value)
		return
	}

	log.Announcef("%s | found %d matching servers - %s", recipe.Name, len(servers), strings.Join(serverNames, ", "))

	for _, a := range recipe.Artifacts {
		log.Infof("→ %s → %s", a.Source, a.Destination)
	}

	for _, c := range recipe.Commands {
		log.Infof("→ $ %s", c)
	}

	answer := log.Questionf("Are you sure you want to continue [y/n]?")
	if utils.ArrayContains([]string{"y", "yes"}, answer, false) == -1 {
		log.Errorf("Quitting, because you said %s", answer)
		os.Exit(0)
	}

	switch recipe.Type {
	case "exec":
		ExecRecipe(servers, recipe, *c.MaxParallelProcesses)
	default:
		log.Errorf("recipe only supports ['exec'] types")
	}
}

func ExecRecipe(servers []config.Server, recipe config.Recipe, maxProcs int) {
	wp := pool.NewPool("ssh")
	wp.AddWorkerGroup(NewSSHWorkerGroup("ssh", 10*time.Millisecond))
	processed := wp.Start(int64(maxProcs))

	startTime := time.Now()
	for _, s := range servers {
		wp.Queue(SSHJob{Server: s, Recipe: recipe})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	count := 0

	for range processed {
		count++
		if count == len(servers) {
			wp.Stop(ctx)
			break
		}
	}

	for {
		if wp.GetWorkerCount() == 0 {
			break
		}
	}

	log.MinorSuccessf("%s | completed in %s", recipe.Name, time.Since(startTime))
}
