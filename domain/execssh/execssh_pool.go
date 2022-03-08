package execssh

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/betas-in/logger"
	"github.com/betas-in/pool"
	"github.com/sudhanshuraheja/golem/domain/artifacts"
	"github.com/sudhanshuraheja/golem/domain/commands"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type SSHPool struct {
	log *logger.CLILogger
}

func NewSSHPool(log *logger.CLILogger) *SSHPool {
	return &SSHPool{
		log: log,
	}
}

func (s *SSHPool) Start(srvs servers.Servers, cmds *[]commands.Command, artfs []*artifacts.Artifact, procs int) {
	log := logger.NewLogger(2, true)

	wp := pool.NewPool("ssh", log)
	wp.AddWorkerGroup(NewSSHWorkerGroup("ssh", s.log, 5*time.Second))

	processed := wp.Start(int64(procs))

	startTime := time.Now()
	for _, s := range srvs {
		wp.Queue(SSHJob{Server: s, Commands: cmds, Artifacts: artfs})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	count := 0

	loop := true
	for loop {
		select {
		case <-processed:
			count++
			if count == len(srvs) {
				wp.Stop(ctx)
				loop = false
				break
			}
		case <-quit:
			wp.Stop(ctx)
			loop = false
			break
		}
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	ticks := 0
	for ; true; <-ticker.C {
		ticks++
		if ticks >= 20 {
			break
		}
		if wp.GetWorkerCount() == 0 {
			break
		}
	}

	s.log.Announce("").Msgf("completed %s", localutils.TimeInSecs(startTime))
}
