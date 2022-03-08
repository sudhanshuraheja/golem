package execssh

import (
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/servers"
)

type SSHJob struct {
	Server    servers.Server
	Commands  *[]commands.Command
	Artifacts []*artifacts.Artifact
}
