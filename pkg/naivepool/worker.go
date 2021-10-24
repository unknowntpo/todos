package naivepool

import (
	"context"
	"sync"
)

// workerChan is the channel connecting between pool and worker, each worker uses 1 workerChan.
// When workerChan is closed, it meas that worker needs to retire.
type workerChan chan jobFunc

type worker struct {
	c workerChan
}

// NewWorker init a new instance of worker.
// we return worker, not *worker to avoid worker escaping to heap, which takes time to do memory-allocation.
func NewWorker(chanSize int) worker {
	return worker{
		c: make(chan jobFunc, chanSize),
	}
}

// worker is the worker that execute the job received from p.workerChan.
func (w *worker) work(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case f := <-w.c:
			f()
		case <-ctx.Done():
			return
		}
	}
	return
}
