package naivepool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkExecute1000Tasks(b *testing.B) {
	maxJobs := 1000
	maxWorkers := 50

	b.Run("naivepool", func(b *testing.B) {
		pool := New(maxJobs, maxWorkers)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pool.Start(ctx)

		var counter uint32
		var wg sync.WaitGroup

		fn := func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}

		b.ResetTimer()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			counter = 0
			for j := 0; j < maxJobs; j++ {
				wg.Add(1)
				pool.Schedule(fn)
			}
			wg.Wait()
		}
		b.StopTimer()
	})
	b.Run("native goroutine", func(b *testing.B) {
		var counter uint32
		var wg sync.WaitGroup

		fn := func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}

		b.ResetTimer()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			counter = 0
			for j := 0; j < maxJobs; j++ {
				wg.Add(1)
				go fn()
			}
			wg.Wait()
		}
		b.StopTimer()
	})
}
