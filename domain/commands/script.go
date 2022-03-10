package commands

import "github.com/sudhanshuraheja/golem/domain/artifacts"

type Script struct {
	Apt      []Apt      `hcl:"apt,block"`
	Commands *[]Command `hcl:"commands"`
	Command  *Command   `hcl:"command"`
}

func (s *Script) Prepare() (Commands, artifacts.Artifacts) {
	_cmds := Commands{}
	_artfs := artifacts.Artifacts{}

	for _, a := range s.Apt {
		cmds, artfs := a.Prepare()
		_cmds.Merge(cmds)
		_artfs.Merge(artfs)
	}

	if s.Commands != nil {
		_cmds.Merge(*s.Commands)
	}

	if s.Command != nil {
		_cmds.Append(*s.Command)
	}

	return _cmds, _artfs
}
