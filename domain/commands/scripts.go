package commands

import "github.com/sudhanshuraheja/golem/domain/artifacts"

type Scripts []*Script

func (s *Scripts) Append(scr Script) {
	*s = append(*s, &scr)
}

func (s *Scripts) Merge(scrs Scripts) {
	*s = append(*s, scrs...)
}

func (s *Scripts) Prepare() (Commands, artifacts.Artifacts) {
	_cmds := Commands{}
	_artfs := artifacts.Artifacts{}

	if s != nil {
		for _, scr := range *s {
			cmds, artfs := scr.Prepare()
			_cmds.Merge(cmds)
			_artfs.Merge(artfs)
		}
	}

	return _cmds, _artfs
}
