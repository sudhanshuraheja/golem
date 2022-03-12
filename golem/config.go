package golem

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

var (
	configBucket = "golemconfig"
	tempdirKey   = "golemconfig.tempdir"
	scriptsKey   = "golemconfig.scripts"
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

	confFile := fmt.Sprintf("%s/global.golem.hcl", golemDir)
	created, err := localutils.Create(confFile)
	if err != nil {
		return err
	}
	if created {
		log.Highlight("golem").Msgf("created %s", confFile)
	}

	return nil
}

func (c *Config) Detect(log *logger.CLILogger, store *kv.Store) ([]string, error) {
	golemDir := c.GolemDir(log)

	files := []string{}

	scriptsPath, err := store.Get(scriptsKey)
	if err != nil {
		return files, err
	}
	dirs := []string{golemDir, scriptsPath, "."}

	for _, dir := range dirs {
		paths, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Error("golem").Msgf("could not read directory <%s>: %v", dir, err)
			continue
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
	return fmt.Sprintf("%s/.config/golem", dirname)
}

func (c *Config) SetupKV(store *kv.Store) error {
	if err := store.Bucket(configBucket); err != nil {
		return err
	}

	tempdir, err := store.Get(tempdirKey)
	if err != nil {
		return err
	}
	if tempdir == "" {
		tempdir, err = ioutil.TempDir("", "golem")
		if err != nil {
			return err
		}
		if err := store.Set(tempdirKey, tempdir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Update(log *logger.CLILogger, store *kv.Store) {
	path := "https://github.com/sudhanshuraheja/golem-scripts/tarball/main"
	downloadPath, err := localutils.Download(log, "config", path)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}

	tempdir, err := store.Get(tempdirKey)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}

	err = utils.Zip().UntarGunzip(downloadPath, tempdir)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}

	folder, err := localutils.FolderNameFromZip(downloadPath)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}

	folderFullpath := fmt.Sprintf("%s/%s", tempdir, folder)
	scriptsPath := fmt.Sprintf("%s/scripts", tempdir)
	err = os.RemoveAll(scriptsPath)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}
	err = os.Rename(folderFullpath, scriptsPath)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}

	err = store.Set(scriptsKey, scriptsPath)
	if err != nil {
		log.Fatal("").Msgf("%v", err)
		os.Exit(1)
	}

	log.Success("update").Msgf("Updated golem scripts in %s", scriptsPath)
}
