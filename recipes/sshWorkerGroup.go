package recipes

import (
	"context"
	"time"

	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/pool"
)

// // WorkerGroup ...
// type SSHWorkerGroup interface {
// 	Process(ctx context.Context, workerCtx *pool.WorkerContext, id string)
// }

type sshWorkerGroup struct {
	name      string
	heartbeat time.Duration
}

// NewWorkerGroup ...
func NewSSHWorkerGroup(name string, heartbeat time.Duration) *sshWorkerGroup {
	w := sshWorkerGroup{
		name:      name,
		heartbeat: heartbeat,
	}
	return &w
}

func (w *sshWorkerGroup) Process(ctx context.Context, workerCtx *pool.WorkerContext, id string) {
	workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Ping: true}
	log.Infof("%s-%s | Started", w.name, id)

	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case j := <-workerCtx.Jobs:

			job, ok := j.(SSHJob)
			if !ok {
				log.Errorf("%s-%s | invalid job", w.name, id)
			}
			// log.Infof("%s-%s | Job %+v", w.name, id, j)

			w.ExecRecipeOnServer(job.Server, job.Recipe)

			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Processed: 1}
			workerCtx.Processed <- j
		case <-ctx.Done():
			log.Successf("%s-%s | Done", w.name, id)
			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Closed: true}
			return
		case <-workerCtx.Close:
			log.Successf("%s-%s | Closing", w.name, id)
			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Closed: true}
			return
		case <-ticker.C:
			log.Tracef("%s-%s | Heartbeat", w.name, id)
			workerCtx.Heartbeat <- pool.Heartbeat{ID: id, Ping: true}
		}
	}
}

func (w *sshWorkerGroup) ExecRecipeOnServer(s config.Server, recipe config.Recipe) {
	ss := SSH{}
	err := ss.Connect(&s)
	if err != nil {
		log.Errorf("%s | %v", s.Name, err)
		return
	}
	ss.Upload(recipe.Artifacts)
	ss.Run(recipe.Commands)
	ss.Close()
}
