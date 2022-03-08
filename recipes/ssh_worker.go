package recipes

import (
	"context"
	"fmt"
	"time"

	"github.com/betas-in/logger"
	"github.com/betas-in/pool"
	"github.com/sudhanshuraheja/golem/domain/servers"
)

type SSHWorkerGroup struct {
	name      string
	log       *logger.CLILogger
	heartbeat time.Duration
}

// NewWorkerGroup ...
func NewSSHWorkerGroup(name string, log *logger.CLILogger, heartbeat time.Duration) *SSHWorkerGroup {
	w := SSHWorkerGroup{
		name:      name,
		log:       log,
		heartbeat: heartbeat,
	}
	return &w
}

func (w *SSHWorkerGroup) Process(ctx context.Context, workerCtx *pool.WorkerContext, id string) {
	workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Ping: true}
	w.log.Trace(w.Name(id)).Msgf("Started")

	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case j := <-workerCtx.Jobs:

			job, ok := j.(SSHJob)
			if !ok {
				w.log.Error(w.Name(id)).Msgf("invalid job")
			}
			w.ExecRecipeOnServer(job.Server, job.Recipe)

			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Processed: 1}
			workerCtx.Processed <- j
		case <-ctx.Done():
			w.log.Trace(w.Name(id)).Msgf("Done")
			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Closed: true}
			return
		case <-workerCtx.Close:
			w.log.Debug(w.Name(id)).Msgf("Closing")
			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Closed: true}
			return
		case <-ticker.C:
			w.log.Trace(w.Name(id)).Msgf("Heartbeat")
			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Ping: true}
		}
	}
}

func (w *SSHWorkerGroup) Name(id string) string {
	return fmt.Sprintf("%s-%s", w.name, id)
}

func (w *SSHWorkerGroup) ExecRecipeOnServer(s servers.Server, recipe *Recipe) {
	ss := SSH{log: w.log}
	err := ss.Connect(&s)
	if err != nil {
		w.log.Error(s.Name).Msgf("%v, please try", err)
		w.log.Success(s.Name).Msgf("$ ssh-keyscan -p %d %s >> ~/.ssh/known_hosts", s.Port, s.PublicIP)
		return
	}
	ss.Upload(recipe.artfs)
	cmds := []string{}
	for _, cmd := range recipe.cmds {
		if cmd.Exec == "" {
			continue
		}
		cmds = append(cmds, cmd.Exec)
	}
	ss.Run(cmds)
	ss.Close()
}
