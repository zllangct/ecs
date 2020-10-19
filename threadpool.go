package ecs

import (
	"math/rand"
	runtime2 "runtime"
	"sync"
)

var globalPool *Pool

//Pool is goroutine pool config.
type Pool struct {
	numWorkers  int32
	jobQueueLen int32
	jobPool     *sync.Pool
	jobQueue    chan *Job
	workerQueue []*Worker
	runtime     *Runtime
}

//get the singleton pool
func GetGlobalPool(runtime *Runtime, numWorkers int, jobQueueLen int) *Pool {
	if globalPool == nil {
		globalPool = NewPool(runtime, numWorkers, jobQueueLen)
	}
	return globalPool
}

//NewPool news goroutine pool
func NewPool(runtime *Runtime, numWorkers int, jobQueueLen int) *Pool {
	if numWorkers == 0 {
		numWorkers = 2 * runtime2.NumCPU()
	}
	if jobQueueLen == 0 {
		jobQueueLen = 20
	}
	jobQueue := make(chan *Job, jobQueueLen)
	workerQueue := make([]*Worker, numWorkers)

	pool := &Pool{
		runtime:     runtime,
		numWorkers:  int32(numWorkers),
		jobQueueLen: int32(jobQueueLen),
		jobQueue:    jobQueue,
		workerQueue: workerQueue,
		jobPool:     &sync.Pool{New: func() interface{} { return &Job{WorkerID: int32(-1)} }},
	}
	pool.Start()
	return pool
}

//random worker, task will run in a random worker
func (p *Pool) AddJob(handler func(*JobContext, ...interface{}), args ...interface{}) {
	job := p.jobPool.Get().(*Job)
	job.Job = handler
	job.Args = args
	job.WorkerID = WORKER_ID_RANDOM
	p.jobQueue <- job
}

//random worker, task will run in a random worker and record the worker id
func (p *Pool) AddJob2(handler func(*JobContext, ...interface{}), args ...interface{}) {
	job := p.jobPool.Get().(*Job)
	job.Job = handler
	job.Args = args

	job.WorkerID = rand.Int31() % p.numWorkers
	p.workerQueue[job.WorkerID].jobQueue <- job
}

//fixed worker,task with the same worker id will push into the same goroutine
func (p *Pool) AddJobFixed(handler func(*JobContext, ...interface{}), args []interface{}, wid int32) {
	job := p.jobPool.Get().(*Job)
	job.Job = handler
	job.Args = args

	if wid <= WORKER_ID_RANDOM || wid >= p.numWorkers {
		job.WorkerID = rand.Int31() % p.numWorkers
		p.workerQueue[job.WorkerID].jobQueue <- job
	} else {
		job.WorkerID = wid
	}
	p.workerQueue[job.WorkerID].jobQueue <- job
}

//Start starts all workers
func (p *Pool) Start() {
	for i := 0; i < cap(p.workerQueue); i++ {
		worker := &Worker{
			runtime:  p.runtime,
			id:       int32(i),
			p:        p,
			jobQueue: make(chan *Job, 10),
			stop:     make(chan struct{}),
		}
		p.workerQueue[i] = worker
		worker.Start()
	}
}

//get the pool size
func (p *Pool) Size() int32 {
	return p.numWorkers
}

//Release release all workers
func (p *Pool) Release() {
	for _, worker := range p.workerQueue {
		worker.stop <- struct{}{}
	}
}
