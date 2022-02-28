package kitchen

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/betas-in/logger"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/recipes"
)

const (
	version = "v0.1.0"
)

type Kitchen struct {
	conf        *config.Config
	configFiles []string
}

func NewKitchen() *Kitchen {
	k := Kitchen{}
	k.conf = &config.Config{}

	err := k.detectConfigFiles()
	if err != nil {
		logger.Fatalf("kitchen | %v", err)
		os.Exit(1)
	}

	if len(k.configFiles) == 0 {
		err := k.initConfigFile()
		if err != nil {
			logger.Fatalf("kitchen | %v", err)
			os.Exit(1)
		}
	}

	for _, path := range k.configFiles {
		conf, err := config.NewConfig(path)
		if err != nil {
			logger.Errorf("%v", err)
		}
		k.mergeConfig(conf)
	}

	k.conf.ResolveServerProvider()
	return &k
}

func (k *Kitchen) mergeConfig(conf *config.Config) {
	k.conf.ServerProviders = append(k.conf.ServerProviders, conf.ServerProviders...)
	k.conf.Servers = append(k.conf.Servers, conf.Servers...)
	k.conf.Recipes = append(k.conf.Recipes, conf.Recipes...)

	if conf.LogLevel != nil {
		k.conf.LogLevel = conf.LogLevel
	}

	if conf.MaxParallelProcesses != nil {
		k.conf.MaxParallelProcesses = conf.MaxParallelProcesses
	}
}

func (k *Kitchen) detectConfigFiles() error {
	files := []string{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find user's home directory: %v", err)
	}
	golemDir := fmt.Sprintf("%s/.golem", homeDir)

	dirs := []string{golemDir, "."}

	for _, dir := range dirs {
		paths, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("could not read directory <%s>: %v", dir, err)
		}

		for _, path := range paths {
			if !path.IsDir() && strings.HasSuffix(path.Name(), ".golem.hcl") {
				fullPath := fmt.Sprintf("%s/%s", dir, path.Name())
				files = append(files, fullPath)
			}
		}
	}

	k.configFiles = files
	return nil
}

func (k *Kitchen) initConfigFile() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find user's home directory: %v", err)
	}
	confDir := fmt.Sprintf("%s/.golem", dirname)

	err = os.MkdirAll(confDir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("could not create conf dir <%s>: %v", confDir, err)
	}

	confFile := fmt.Sprintf("%s/config.golem.hcl", confDir)
	_, err = os.Stat(confFile)
	if os.IsNotExist(err) {
		file, err := os.Create(confFile)
		if err != nil {
			return fmt.Errorf("error creating conf file <%s>: %v", confFile, err)
		}
		defer file.Close()
		logger.MinorSuccessf("kitchen | conf file created at %s", confFile)
	} else if err != nil {
		return fmt.Errorf("error checking conf file <%s>: %v", confFile, err)
	}
	return nil
}

func (k *Kitchen) Exec(recipe string) {
	if recipe != "" && k.conf != nil && k.conf.MaxParallelProcesses != nil {
		logger.Announcef("%s | running recipe with max %d routines", recipe, *k.conf.MaxParallelProcesses)
	}
	r := recipes.NewRecipes(k.conf)
	switch recipe {
	case "":
		// log.MinorSuccessf("We found these recipes in '~/.golem' and '.'")
		r.Servers()
		r.List()
	case "version":
		logger.MinorSuccessf("golem version: %s", version)
	case "list":
		r.List()
	case "servers":
		r.Servers()
	default:
		r.Run(recipe)
	}
}
