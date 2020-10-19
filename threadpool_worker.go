package ecs

const WORKER_ID_RANDOM int32 = -1

//Worker goroutine struct.
type Worker struct {
	runtime  *Runtime
	id       int32
	p        *Pool
	jobQueue chan *Job
	stop     chan struct{}
}

//Start start goroutine pool.
func (w *Worker) Start() {
	go func() {
		var job *Job
		for {
			select {
			case job = <-w.jobQueue:
			case job = <-w.p.jobQueue:
				//task which worker id not nil will push into the target goroutine to insure data safety
				if job.WorkerID != WORKER_ID_RANDOM {
					if job.WorkerID >= 0 && job.WorkerID < w.p.numWorkers {
						w.p.workerQueue[job.WorkerID].jobQueue <- job
						continue
					}
				}
			case <-w.stop:
				return
			}
			ctx := &JobContext{
				WorkerID: job.WorkerID,
				Runtime:  w.runtime,
			}
			if w.runtime.config.Debug {
				job.Job(ctx, job.Args...)
			} else {
				err := Try(func() {
					job.Job(ctx, job.Args...)
				})
				if err != nil && w.runtime.logger != nil {
					w.runtime.logger.Error(err)
				}
			}
			w.p.jobPool.Put(job.Init())
		}
	}()
}
