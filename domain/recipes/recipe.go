package recipes

import (
	"fmt"
	"os"
	"strings"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/execcmd"
	"github.com/sudhanshuraheja/golem/domain/execssh"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/template"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

var (
	tiny = 100
)

type Recipe struct {
	Name       string                `hcl:"name,label"`
	Type       string                `hcl:"type,label"`
	Match      *servers.Match        `hcl:"match,block"`
	KeyValues  []*kv.KeyValue        `hcl:"kv,block"`
	Scripts    []*commands.Script    `hcl:"script,block"`
	Artifacts  []*artifacts.Artifact `hcl:"artifact,block"`
	Commands   *[]commands.Command   `hcl:"commands"`
	SourceFile string
}

func (r *Recipe) Prepare(log *logger.CLILogger, store *kv.Store) error {
	_cmds := commands.Commands{}
	_artfs := artifacts.Artifacts{}

	// Scripts
	for _, s := range r.Scripts {
		cmds, artfs := s.Prepare()
		_cmds.Merge(cmds)
		_artfs.Merge(artfs)
	}

	// Artifacts
	for _, a := range r.Artifacts {
		_artfs.Append(*a)
	}

	// Commands
	if r.Commands != nil {
		for _, c := range *r.Commands {
			_cmds.Append(c)
		}
	}

	cmds := _cmds.ToArray()
	r.Commands = &cmds
	r.Artifacts = _artfs

	return nil
}

func (r *Recipe) PrepareForExecution(log *logger.CLILogger, tpl *template.Template, store *kv.Store) error {
	_cmds := commands.Commands{}
	_artfs := artifacts.Artifacts{}

	// KeyValues
	for _, k := range r.KeyValues {
		setup, err := k.PrepareForExecution(store)
		if err != nil {
			log.Error(r.Name).Msgf("%v", err)
			os.Exit(1)
		}
		if setup {
			log.Info(r.Name).Msgf("setup key %s in store", k.Path)
		}
	}

	// Artifacts
	for _, a := range r.Artifacts {
		err := a.PrepareForExecution(log, tpl)
		if err != nil {
			log.Error(r.Name).Msgf("%v", err)
		}
		_artfs.Append(*a)
	}

	// Commands
	if r.Commands != nil {
		for _, c := range *r.Commands {
			cmd, err := c.PrepareForExecution(tpl)
			if err != nil {
				return err
			}
			_cmds.Append(*cmd)
		}
	}

	cmds := _cmds.ToArray()
	r.Commands = &cmds
	r.Artifacts = _artfs.ToPointerArray()

	return nil
}

func (r *Recipe) Execute(log *logger.CLILogger, srvs servers.Servers, procs int) {
	switch r.Type {
	case "remote":
		pool := execssh.NewSSHPool(log)
		pool.Start(srvs, r.Commands, r.Artifacts, procs)
	case "local":
		pool := execcmd.NewExecCmd(log)
		pool.Start(r.Commands, r.Artifacts)
	default:
		log.Error(r.Name).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipe) Display(log *logger.CLILogger, tpl *template.Template, query string) {
	if query != "" && !strings.Contains(r.Name, query) {
		return
	}

	match := ""
	if r.Match != nil {
		match = fmt.Sprintf("%s %s %s", r.Match.Attribute, r.Match.Operator, r.Match.Value)
		match = logger.Yellow(match)
	}

	log.Info("recipe").Msgf("%s %s %s", logger.CyanBold(r.Name), r.Type, match)
	log.Info("").Msgf("%s %s", logger.Cyan("source"), r.SourceFile)

	for _, keyvalue := range r.KeyValues {
		if keyvalue != nil {
			log.Info("").Msgf(
				"%s %s: %s",
				logger.Yellow("kv"),
				logger.CyanBold(keyvalue.Path),
				keyvalue.Value,
			)
		}
	}

	if r.Artifacts != nil {
		for _, artf := range r.Artifacts {
			source := artf.GetSource()

			log.Info("").Msgf(
				"%s %s %s %s",
				logger.Cyan("uploading"),
				localutils.TinyString(source, tiny),
				logger.Cyan("to"),
				localutils.TinyString(artf.Destination, tiny),
			)
		}
	}

	if r.Commands != nil {
		for _, command := range *r.Commands {
			exec, err := tpl.Execute(string(command))
			if err != nil {
				log.Error(r.Name).Msgf("could not parse template %s: %v", command, err)
			}
			exec = strings.TrimSuffix(exec, "\n")
			log.Info("").Msgf("%s %s", logger.Cyan("$"), exec)
		}
	}
}

func (r *Recipe) AskPermission(log *logger.CLILogger) {
	answer := localutils.Question(log, "", "Are you sure you want to continue [y/n]?")
	if utils.Array().Contains([]string{"y", "yes"}, answer, false) == -1 {
		log.Error("").Msgf("Quitting, because you said %s", answer)
		if answer != "EOF" {
			os.Exit(0)
		}
	}
}
