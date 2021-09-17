package ecs

import (
	runtime2 "runtime"
)

//Worker goroutine struct.
type Worker struct {
	p        *Pool
	jobQueue chan func()
	stop     chan struct{}
}

//Start goroutine pool.
func (w *Worker) Start() {
	go func() {
		var job func()
		for {
			select {
			case job = <-w.jobQueue:
			case job = <-w.p.jobQueue:
			case <-w.stop:
				return
			}
			Try(job)
		}
	}()
}

//Pool is goroutine pool config.
type Pool struct {
	size         uint32
	jobQueueSize uint32
	jobQueue     chan func()
	workers      []*Worker
}

//NewPool news goroutine pool
func NewPool(size uint32, jobQueueSize uint32) *Pool {
	if size == 0 {
		size = uint32(2 * runtime2.NumCPU())
	}
	if jobQueueSize == 0 {
		jobQueueSize = 20
	}
	jobQueue := make(chan func(), jobQueueSize)
	workerQueue := make([]*Worker, size)

	pool := &Pool{
		size:         uint32(size),
		jobQueueSize: uint32(jobQueueSize),
		jobQueue:     jobQueue,
		workers:      workerQueue,
	}
	pool.Start()
	return pool
}

// Add hashKey is an optional parameter, job will be executed in a random worker
// when hashKey is regardless, in fixed worker calculated by hash when hashKey is
// specified
func (p *Pool) Add(job func(), hashKey ...uint32) {
	if len(hashKey) > 0 {
		p.workers[hashKey[0]%p.size].jobQueue <- job
		return
	}
	p.jobQueue <- job
}

//Start all workers
func (p *Pool) Start() {
	for i := 0; i < cap(p.workers); i++ {
		worker := &Worker{
			p:        p,
			jobQueue: make(chan func(), p.jobQueueSize),
			stop:     make(chan struct{}),
		}
		p.workers[i] = worker
		worker.Start()
	}
}

// Size get the pool size
func (p *Pool) Size() uint32 {
	return p.size
}

//Release stop all workers
func (p *Pool) Release() {
	for _, worker := range p.workers {
		worker.stop <- struct{}{}
	}
}
