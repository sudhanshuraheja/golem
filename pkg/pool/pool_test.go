package pool

import (
	"context"
	"testing"
	"time"

	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func TestWorkerPool(t *testing.T) {
	queueSize := 200
	wp := NewPool("wk")
	wp.AddWorkerGroup(NewWorkerGroup("wk", time.Second))
	processed := wp.Start(2)

	wp.Update(5)
	for i := 0; i < queueSize; i++ {
		wp.Queue(100 + i)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count := 0
	for range processed {
		count++
		if count == queueSize {
			utils.Equals(t, queueSize, count)
			wp.Stop(ctx)
			break
		}
	}

}
