package ecs

//job type: Parallel;Serial
type JobType int

const JOB_TYPE_PARALLEL JobType = 1
const JOB_TYPE_SERIAL JobType = 2
const JOB_TYPE_DEFAULT = JOB_TYPE_PARALLEL

//Job is a function for doing jobs.
type Job struct {
	WorkerID int32
	Args     []interface{}
	Job      func(ctx []interface{}, args ...interface{})
}

//initialize the job
func (p *Job) Init() *Job {
	p.WorkerID = -1
	p.Args = nil
	p.Job = nil
	return p
}
