package pool

import (
	"context"
	"log"
	"time"
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
	log.Printf("[pool.wg.process][%s-%s] Started", w.name, id)

	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case j := <-workerCtx.Jobs:
			log.Printf("[pool.wg.process][%s-%s] Job %+v", w.name, id, j)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Processed: 1}
			workerCtx.Processed <- j
		case <-ctx.Done():
			log.Printf("[pool.wg.process][%s-%s] Done", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-workerCtx.Close:
			log.Printf("pool.wg.process[%s-%s] Closing", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Closed: true}
			return
		case <-ticker.C:
			log.Printf("pool.wg.process[%s-%s] Heartbeat", w.name, id)
			workerCtx.Heartbeat <- Heartbeat{ID: id, Ping: true}
		}
	}
}
