package ecs

type JobContext struct {
	WorkerID int32
	Runtime  *Runtime
}
