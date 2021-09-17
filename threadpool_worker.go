package ecs


const WORKER_ID_RANDOM int32 = -1

const JOB_TYPE_PARALLEL JobType = 1
const JOB_TYPE_SERIAL JobType = 2
const JOB_TYPE_DEFAULT = JOB_TYPE_PARALLEL

// JobType job type: Parallel;Serial
type JobType int

type JobContext struct {
	WorkerID int32
	Runtime  *ecsRuntime
}

//Job is a function for doing jobs.
type Job struct {
	WorkerID int32
	Args     []interface{}
	Job      func(ctx JobContext, args ...interface{})
}

// Init initialize the job
func (p *Job) Init() *Job {
	p.WorkerID = -1
	p.Args = nil
	p.Job = nil
	return p
}

//Worker goroutine struct.
type Worker struct {
	runtime  *ecsRuntime
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
			ctx := JobContext{
				WorkerID: job.WorkerID,
				Runtime:  w.runtime,
			}
			err := Try(func() {
				job.Job(ctx, job.Args...)
			})
			if err != nil && w.runtime.logger != nil {
				w.runtime.logger.Error(err)
			}
			w.p.jobPool.Put(job.Init())
		}
	}()
}
