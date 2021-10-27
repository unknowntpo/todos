package naivepool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

func Example_adder() {
	maxJobs := 1000
	maxWorkers := 50
	workerChanSize := 10

	pool := New(maxJobs, maxWorkers, workerChanSize)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	var counter uint32
	var wg sync.WaitGroup

	fn := func() {
		defer wg.Done()
		atomic.AddUint32(&counter, 1)
	}

	counter = 0
	for j := 0; j < maxJobs; j++ {
		wg.Add(1)
		pool.Schedule(fn)
	}
	wg.Wait()

	// Call cancel to stop the pool explicitly.
	cancel()

	// Wait for all workers to retire.
	pool.wg.Wait()

	fmt.Println(counter)
	// Output: 1000
}
