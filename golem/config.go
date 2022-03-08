package golem

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/betas-in/logger"
)

type Config struct {
	Recipe string `arg:"positional"`
	Param1 string `arg:"positional"`
	Param2 string `arg:"positional"`
	Param3 string `arg:"positional"`
	Param4 string `arg:"positional"`
}

func (c *Config) Init(log *logger.CLILogger) error {
	golemDir := c.GolemDir(log)

	err := os.MkdirAll(golemDir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("could not create conf dir <%s>: %v", golemDir, err)
	}

	confFile := fmt.Sprintf("%s/config.golem.hcl", golemDir)
	_, err = os.Stat(confFile)
	if os.IsNotExist(err) {
		file, err := os.Create(confFile)
		if err != nil {
			return fmt.Errorf("error creating conf file <%s>: %v", confFile, err)
		}
		defer file.Close()
		log.Highlight("golem").Msgf("conf file created at %s", confFile)
	} else if err != nil {
		return fmt.Errorf("error checking conf file <%s>: %v", confFile, err)
	}
	return nil
}

func (c *Config) Detect(log *logger.CLILogger) ([]string, error) {
	golemDir := c.GolemDir(log)
	dirs := []string{golemDir, "."}

	files := []string{}
	for _, dir := range dirs {
		paths, err := ioutil.ReadDir(dir)
		if err != nil {
			return files, fmt.Errorf("could not read directory <%s>: %v", dir, err)
		}

		for _, path := range paths {
			if !path.IsDir() && strings.HasSuffix(path.Name(), ".golem.hcl") {
				fullPath := fmt.Sprintf("%s/%s", dir, path.Name())
				files = append(files, fullPath)
			}
		}
	}

	return files, nil
}

func (c *Config) GolemDir(log *logger.CLILogger) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("").Msgf("could not find user's home directory: %v", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%s/.golem", dirname)
}
