package background

import (
	"context"
	"sync"
)

type taskFunc func()

type Background struct {
	wg       sync.WaitGroup
	taskChan chan taskFunc
}

// Start starts handling background tasks.
func (b *Background) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case task := <-b.taskChan:
				go func() {
					// use context to stop go routine with 3 seconds' deadline?
					b.wg.Add(1)
					task()
					b.wg.Done()
				}()
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
}

// Wait waits for all tasks to be completed.
// FIXME: How to wait until deadline comes ?
func (b *Background) Wait() {
	b.wg.Wait()
}

// Add adds new background task.
func (b *Background) Add(fn taskFunc) {
	b.taskChan <- fn
}
