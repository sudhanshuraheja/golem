package pool

import (
	"context"
	"time"

	"github.com/sudhanshuraheja/golem/pkg/log"
)

// WorkerGroup ...
type WorkerGroup interface {
	Process(ctx context.Context, workerCtx *WorkerContext, id string)
}

type workerGroup struct {
	name      string
	heartbeat time.Duration
}

// NewWorkerGroup ...
func NewWorkerGroup(name string, heartbeat time.Duration) WorkerGroup {
	w := workerGroup{
		name:      name,
		heartbeat: heartbeat,
	}
	return &w
}

func (w *workerGroup) Process(ctx context.Context, workerCtx *WorkerContext, id string) {
	workerCtx.Heartbeat <- Heartbeat{ID: id, Ping: true}
	log.Infof("pool | %s-%s | Started", w.name, id)

	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case j := <-workerCtx.Jobs:
			log.Infof("pool | %s-%s | Job %+v", w.name, id, j)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Processed: 1}
			workerCtx.Processed <- j
		case <-ctx.Done():
			log.Successf("pool | %s-%s | Done", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-workerCtx.Close:
			log.Successf("pool | %s-%s | Closing", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-ticker.C:
			log.Tracef("pool | %s-%s | Heartbeat", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Ping: true}
		}
	}
}
