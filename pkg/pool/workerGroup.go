package pool

import (
	"context"
	"time"

	"github.com/sudhanshuraheja/golem/pkg/logger"
)

// WorkerGroup ...
type WorkerGroup interface {
	Process(ctx context.Context, workerCtx *WorkerContext, id string)
}

type workerGroup struct {
	log       *logger.Logger
	name      string
	heartbeat time.Duration
}

// NewWorkerGroup ...
func NewWorkerGroup(name string, heartbeat time.Duration, log *logger.Logger) WorkerGroup {
	w := workerGroup{
		name:      name,
		heartbeat: heartbeat,
		log:       log,
	}
	return &w
}

func (w *workerGroup) Process(ctx context.Context, workerCtx *WorkerContext, id string) {
	workerCtx.Heartbeat <- Heartbeat{ID: id, Ping: true}
	w.log.Info("pool.wg.process").Msgf("[%s-%s] Started", w.name, id)

	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case j := <-workerCtx.Jobs:
			w.log.Info("pool.wg.process").Msgf("[%s-%s] Job %+v", w.name, id, j)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Processed: 1}
			workerCtx.Processed <- j
		case <-ctx.Done():
			w.log.Info("pool.wg.process").Msgf("[%s-%s] Done", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-workerCtx.Close:
			w.log.Info("pool.wg.process").Msgf("[%s-%s] Closing", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-ticker.C:
			w.log.Info("pool.wg.process").Msgf("[%s-%s] Heartbeat", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Ping: true}
		}
	}
}
