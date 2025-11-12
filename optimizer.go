package ecs

import (
	"fmt"
	"sort"
	"time"
)

type ShapeInfo struct {
	execCount int64
	shape     Compound
}

type optReporter struct {
	shapeUsage map[FixedCompound]int64
}

func newOptReporter() *optReporter {
	return &optReporter{shapeUsage: make(map[FixedCompound]int64)}
}

func (o *optReporter) shapeUsageAdd(compound Compound) {
	key := NewFixedCompound(compound)
	if _, ok := o.shapeUsage[key]; ok {
		o.shapeUsage[key]++
	} else {
		o.shapeUsage[key] = 1
	}
}

type optimizer struct {
	world                  *world
	startTime              time.Time
	expireTime             time.Time
	lastSample             time.Time
	shapeInfos             []ShapeInfo
	lastCollectConsumption time.Duration
}

func newOptimizer(world *world) *optimizer {
	return &optimizer{world: world}
}

func (o *optimizer) collect() {
	start := time.Now()
	var opts []*optReporter
	for _, system := range o.world.systems.systems {
		if system != nil {
			opts = append(opts, system.getOptReporter())
		}
	}
	//all shapes
	var shapeRef = map[FixedCompound]int64{}
	for _, opt := range opts {
		for shp, count := range opt.shapeUsage {
			shapeRef[shp] += count
		}
	}
	//sort
	o.shapeInfos = make([]ShapeInfo, len(shapeRef))
	for compound, count := range shapeRef {
		o.shapeInfos = append(o.shapeInfos, ShapeInfo{
			execCount: count,
			shape:     compound.Compound(),
		})
	}
	sort.Slice(o.shapeInfos, func(i, j int) bool {
		return o.shapeInfos[i].execCount > o.shapeInfos[j].execCount
	})

	o.lastCollectConsumption = time.Since(start)
}

func (o *optimizer) optimize(IdleTime time.Duration, force bool) {
	fmt.Printf("start optimize, rest time: %v\n", IdleTime)
	now := time.Now()
	o.startTime = now
	o.expireTime = o.startTime.Add(IdleTime)
	elapsed := o.elapsedStep()

	o.collect()
	rest := o.expire()
	elapsed = o.elapsedStep()
	fmt.Printf("collect step: %v, rest time: %v\n", elapsed, rest)

	o.memTidy(force)
	elapsed = o.elapsedStep()

	rest = o.expire()
	total := time.Now().Sub(o.startTime)
	fmt.Printf("end optimize, rest time: %v, total: %v\n", rest, total)
}

func (o *optimizer) expire() time.Duration {
	return time.Until(o.expireTime)
}

func (o *optimizer) elapsed() time.Duration {
	return time.Now().Sub(o.startTime)
}

func (o *optimizer) elapsedStep() time.Duration {
	now := time.Now()
	r := now.Sub(o.lastSample)
	o.lastSample = now
	return r
}

func (o *optimizer) memTidy(force bool) {
	rest := o.expire()
	if !force && rest < time.Millisecond {
		return
	}

	for _, collection := range o.world.components {
		collection.Sort()
		if !force && o.expire() < time.Millisecond {
			break
		}
	}
}
