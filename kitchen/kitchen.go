package kitchen

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/betas-in/logger"

	"github.com/sudhanshuraheja/golem/config"
)

type Kitchen struct {
	cliConf     *CLIConfig
	conf        *config.Config
	log         *logger.CLILogger
	configFiles []string
}

type CLIConfig struct {
	Recipe string `arg:"positional"`
	Param1 string `arg:"positional"`
	Param2 string `arg:"positional"`
	Param3 string `arg:"positional"`
	Param4 string `arg:"positional"`
}

func NewKitchen(cliConf *CLIConfig) {
	k := Kitchen{}
	k.cliConf = cliConf
	k.log = logger.NewCLILogger(6, 12)
	k.conf = &config.Config{}

	err := k.detectConfigFiles()
	if err != nil {
		k.log.Fatal("golem").Msgf("%v", err)
		os.Exit(1)
	}

	if len(k.configFiles) == 0 {
		err := k.initConfigFile()
		if err != nil {
			k.log.Fatal("golem").Msgf("%v", err)
			os.Exit(1)
		}
	}

	for _, path := range k.configFiles {
		conf := config.NewConfig(path)
		k.mergeConfig(conf)
	}

	err = k.conf.ResolveServerProviders()
	if err != nil {
		k.log.Error("golem").Msgf("%v", err)
	}

	if k.conf.LogLevel != nil {
		k.log = logger.NewCLILogger(*k.conf.LogLevel, 12)
	}
	k.Exec()
}

func (k *Kitchen) Exec() {
	r := NewRecipes(k.conf, k.log)
	switch k.cliConf.Recipe {
	case "":
		r.ListRecipes(k.cliConf.Param1)
	case "version":
		k.log.Highlight("golem").Msgf("version: %s", version)
	case "list":
		r.ListRecipes(k.cliConf.Param1)
	case "servers":
		r.ListServers(k.cliConf.Param1)
	case "kv":
		r.KV(k.cliConf.Param1, k.cliConf.Param2)
	default:
		if k.cliConf.Recipe != "" && k.conf != nil && k.conf.MaxParallelProcesses != nil {
			k.log.Announce(k.cliConf.Recipe).Msgf("running with a maximum of %d routines %s", *k.conf.MaxParallelProcesses, logger.CyanBold(k.cliConf.Recipe))
		}
		r.Run(k.cliConf.Recipe)
	}
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

	if conf.Vars != nil {
		if k.conf.Vars == nil {
			vars := make(map[string]string)
			k.conf.Vars = &vars
		}
		for key, value := range *conf.Vars {
			(*k.conf.Vars)[key] = value
		}
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
		k.log.Highlight("golem").Msgf("conf file created at %s", confFile)
	} else if err != nil {
		return fmt.Errorf("error checking conf file <%s>: %v", confFile, err)
	}
	return nil
}
