package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type jobFunc func()

type Pool struct {
	// we use jobChan to communicate between caller of Pool and Pool.
	jobChan chan jobFunc
	// we use workerChan to send job from Pool to workers.
	workerChan chan jobFunc
	maxJobs    int
	maxWorkers int
	// use waitgroup to wait for workers done its job and retire.
	wg sync.WaitGroup
}

// Init inits goroutine pool with capacity of jobchan and workerchan
func New(maxJobs, maxWorkers int) *Pool {
	p := &Pool{
		jobChan:    make(chan jobFunc, maxJobs),
		workerChan: make(chan jobFunc, maxWorkers),
		maxJobs:    maxJobs,
		maxWorkers: maxWorkers,
	}

	return p
}
func (p *Pool) Start(ctx context.Context) {
	// Dispatch workers.
	p.wg.Add(p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		go p.worker(ctx)
	}

	go func() {
		for {
			select {
			// Received a job.
			// Dispatch it to workers.
			case job := <-p.jobChan:
				fmt.Println("pool send a job to workerChan")
				p.workerChan <- job
			case <-ctx.Done():
				fmt.Println("Pool received cancel")
				return
			default:
			}
		}

	}()
	return
}

// Wait waits for all workers finish its job and retire.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Schedule sends the job the p.jobChan.
func (p *Pool) Schedule(job jobFunc) {
	fmt.Println("I send a job to pool")

	p.jobChan <- job
}

// worker is the worker that execute the job received from p.workerChan.
func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case job := <-p.workerChan:
			fmt.Println("worker received a job")
			job()
		case <-ctx.Done():
			fmt.Println("worker retired!")
			return
		}
	}
}
