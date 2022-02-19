package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sudhanshuraheja/golem/pkg/log"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

var (
	procStarting   = "starting"
	procRunning    = "running"
	maxChannelSize = 99999
)

// Pool ...
type Pool interface {
	Queue(data interface{})
	AddWorkerGroup(wg WorkerGroup)
	Start(count int64) chan interface{}
	Stop(ctx context.Context)
	Update(count int64)
	GetWorkerCount() int64
	GetHeartbeat() chan Heartbeat
}

type pool struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	Name         string
	workerCount  int64
	closingCount int64
	workerGroup  WorkerGroup
	workers      map[string]*Worker
	jobs         chan interface{}
	processed    chan interface{}
	close        chan struct{}
	heartbeat    chan Heartbeat
	Heartbeat    chan Heartbeat
	lock         *sync.RWMutex
	listening    bool
}

// Worker ...
type Worker struct {
	ID            string
	State         string
	LastHeartBeat int64
}

// Heartbeat ...
type Heartbeat struct {
	ID        string
	Processed int64
	Closed    bool
	Ping      bool
}

// WorkerContext ...
type WorkerContext struct {
	Heartbeat chan Heartbeat
	Jobs      chan interface{}
	Processed chan interface{}
	Close     chan struct{}
}

// NewPool ...
func NewPool(name string) Pool {
	w := pool{
		Name:         name,
		workerCount:  0,
		closingCount: 0,
		workers:      map[string]*Worker{},
		jobs:         make(chan interface{}, maxChannelSize),
		processed:    make(chan interface{}, maxChannelSize),
		close:        make(chan struct{}, maxChannelSize),
		heartbeat:    make(chan Heartbeat, maxChannelSize),
		lock:         &sync.RWMutex{},
	}
	w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	return &w
}

func (w *pool) Queue(job interface{}) {
	w.jobs <- job
}

func (w *pool) AddWorkerGroup(workerGroup WorkerGroup) {
	w.workerGroup = workerGroup
}

func (w *pool) Start(count int64) chan interface{} {
	go w.listen()
	w.listening = true
	w.Update(count)
	return w.processed
}

func (w *pool) Update(count int64) {
	log.Infof("pool | %s | updating worker count from <%d> to <%d>", w.Name, w.getExpectedWorkerCount(), count)
	w.setExpectedWorkerCount(count)

	for w.getExpectedWorkerCount() != w.getActualWorkerCount() {
		if w.getExpectedWorkerCount() < w.getActualWorkerCount() {
			w.removeWorker()
		} else if w.getExpectedWorkerCount() > w.getActualWorkerCount() {
			w.addWorker()
		}
	}
}

func (w *pool) GetWorkerCount() int64 {
	return w.getActualWorkerCount()
}

func (w *pool) Stop(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond) // #TODO remove hard coding
	defer ticker.Stop()
	w.ctxCancel()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if w.getActualWorkerCount() == 0 && !w.getListening() {
				return
			}
		}
	}
}

func (w *pool) GetHeartbeat() chan Heartbeat {
	w.Heartbeat = make(chan Heartbeat, maxChannelSize)
	return w.Heartbeat
}

//
// Internal Functions
//

func (w *pool) getNewWorker(id, state string) *Worker {
	return &Worker{ID: id, State: state}
}

func (w *pool) addWorker() {
	w.lock.Lock()
	id := w.getUniqueID()
	w.workers[id] = w.getNewWorker(id, procStarting)
	w.lock.Unlock()

	go w.workerGroup.Process(w.ctx, w.getNewWorkerContext(), id)
}

func (w *pool) removeWorker() {
	w.close <- struct{}{}
	w.closingCount++
}

func (w *pool) listen() {
	for hb := range w.heartbeat {
		switch {
		case hb.Ping:
			w.updateHB(hb.ID, procRunning, 0)
		case hb.Closed:
			w.updateHB(hb.ID, procRunning, 0)
			w.lock.Lock()
			if w.closingCount > 0 {
				w.closingCount--
			}
			delete(w.workers, hb.ID)
			w.lock.Unlock()
		case hb.Processed > 0:
			w.updateHB(hb.ID, procRunning, hb.Processed)
		}
		if w.Heartbeat != nil {
			w.Heartbeat <- hb
		}
		if w.getActualWorkerCount() == 0 {
			close(w.heartbeat)
			if w.Heartbeat != nil {
				close(w.Heartbeat)
			}
		}
		w.emptyProcessed()
	}
	w.lock.Lock()
	w.listening = false
	w.lock.Unlock()
}

func (w *pool) emptyProcessed() {
	w.lock.Lock()
	defer w.lock.Unlock()
	if float64(len(w.processed)) > (float64(maxChannelSize) * 0.9) {
		fmt.Println("w.processed", len(w.processed), "max", float64(maxChannelSize)*0.9)

		count := 0
		for float64(len(w.processed)) > (float64(maxChannelSize) * 0.5) {
			<-w.processed
			count++
		}
		log.Infof("pool | %s | w.processed at %d, removed %d", w.Name, len(w.processed), count)
	}
}

func (w *pool) updateHB(id string, state string, processed int64) {
	w.lock.Lock()
	defer w.lock.Unlock()
	worker, ok := w.workers[id]
	if !ok {
		worker = w.getNewWorker(id, state)
	}
	worker.State = state
	worker.LastHeartBeat = time.Now().Unix()
	w.workers[id] = worker
}

func (w *pool) getExpectedWorkerCount() int64 {
	return atomic.LoadInt64(&w.workerCount)
}

func (w *pool) setExpectedWorkerCount(n int64) {
	atomic.StoreInt64(&w.workerCount, n)
}

func (w *pool) getActualWorkerCount() int64 {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return int64(len(w.workers)) - w.closingCount
}

func (w *pool) getListening() bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.listening
}

func (w *pool) getUniqueID() string {
	var id string
	for id == "" {
		id = utils.GetShortUUID()
		if w.doesIDExist(id) {
			id = ""
		}
	}
	return id
}

func (w *pool) doesIDExist(id string) bool {
	_, ok := w.workers[id]
	return ok
}

func (w *pool) getNewWorkerContext() *WorkerContext {
	return &WorkerContext{
		Heartbeat: w.heartbeat,
		Jobs:      w.jobs,
		Processed: w.processed,
		Close:     w.close,
	}
}
