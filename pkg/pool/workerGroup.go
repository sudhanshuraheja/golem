package pool

import (
	"context"
	"time"

	"github.com/betas-in/logger"
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
	logger.Infof("%s-%s | Started", w.name, id)

	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case j := <-workerCtx.Jobs:
			logger.Infof("%s-%s | Job %+v", w.name, id, j)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Processed: 1}
			workerCtx.Processed <- j
		case <-ctx.Done():
			logger.Successf("%s-%s | Done", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-workerCtx.Close:
			logger.Successf("%s-%s | Closing", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-ticker.C:
			logger.Tracef("%s-%s | Heartbeat", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Ping: true}
		}
	}
}
