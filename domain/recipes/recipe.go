package recipes

import (
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

type Recipe struct {
	Name      string                `hcl:"name,label"`
	Type      string                `hcl:"type,label"`
	Match     *servers.Match        `hcl:"match,block"`
	KeyValues []*kv.KeyValue        `hcl:"kv,block"`
	Scripts   []*commands.Script    `hcl:"script,block"`
	Artifacts []*artifacts.Artifact `hcl:"artifact,block"`
	Commands  *[]commands.Command   `hcl:"commands"`
}

func (r *Recipe) Prepare(log *logger.CLILogger, store *kv.Store) error {
	_cmds := commands.Commands{}
	_artfs := artifacts.Artifacts{}

	// KeyValues
	for _, k := range r.KeyValues {
		setup, err := k.Setup(store)
		if err != nil {
			return err
		}
		if setup {
			log.Info(r.Name).Msgf("setup key %s in store", k.Path)
		}
	}

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

func (r *Recipe) PrepareForExecution(log *logger.CLILogger, tpl *template.Template) error {
	_cmds := commands.Commands{}
	_artfs := artifacts.Artifacts{}

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
		log.Error(string(r.Name)).Msgf("recipe only supports ['remote', 'local'] types")
	}
}

func (r *Recipe) Display(log *logger.CLILogger, tpl *template.Template, query string) {
	if query != "" && !strings.Contains(string(r.Name), query) {
		return
	}

	log.Info("recipe").Msgf("%s", logger.CyanBold(r.Name))

	if r.Artifacts != nil {
		for _, artf := range r.Artifacts {
			source := artf.GetSource()

			log.Info("").Msgf(
				"%s %s %s %s",
				logger.Cyan("uploading"),
				localutils.TinyString(source, 50),
				logger.Cyan("to"),
				localutils.TinyString(artf.Destination, 50),
			)
		}
	}

	if r.Commands != nil {
		for _, command := range *r.Commands {
			exec, err := tpl.Execute(string(command))
			if err != nil {
				log.Error(r.Name).Msgf("could not parse template %s: %v", command, err)
			}
			log.Info("").Msgf("$ %s", localutils.TinyString(exec, 100))
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
