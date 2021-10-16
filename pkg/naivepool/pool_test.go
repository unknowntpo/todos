package naivepool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkExecute10000Tasks(b *testing.B) {
	maxJobs := 1000
	maxWorkers := 1000

	pool := New(maxJobs, maxWorkers)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)
	b.Run("naivepool", func(b *testing.B) {
		var counter uint32
		var wg sync.WaitGroup

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			counter = 0
			for j := 0; j < maxJobs; j++ {
				wg.Add(1)
				pool.Schedule(func() {
					defer wg.Done()
					atomic.AddUint32(&counter, 1)
				})
			}
			wg.Wait()
		}
	})
	b.Run("native goroutine", func(b *testing.B) {
		var counter uint32
		var wg sync.WaitGroup

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			counter = 0
			for j := 0; j < maxJobs; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					atomic.AddUint32(&counter, 1)
				}()
			}
			wg.Wait()
		}
	})
}
