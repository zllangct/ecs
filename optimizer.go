package ecs

import (
	"reflect"
	"sort"
	"time"
)

type ShapeInfo struct {
	typ    reflect.Type
	eNum   int64
	shapes []IShapeGetter
}

type OptimizerReporter struct {
	shapeUsage map[reflect.Type]IShapeGetter
}

func (o *OptimizerReporter) init() {
	o.shapeUsage = map[reflect.Type]IShapeGetter{}
}

type optimizer struct {
	world                  *ecsWorld
	expireTime             time.Time
	shapeInfos             []*ShapeInfo
	lastCollectConsumption time.Duration
}

func newOptimizer(world *ecsWorld) *optimizer {
	return &optimizer{world: world}
}

// Collect 采集分布在各个系统中的OptimizerReporter
func (o *optimizer) collect() {
	start := time.Now()
	var opts []*OptimizerReporter
	o.world.systemFlow.systems.Range(func(key, value interface{}) bool {
		system, ok := value.(ISystem)
		if !ok {
			return true
		}
		if system != nil {
			opts = append(opts, system.getOptimizer())
		}
		return true
	})
	//all shapes
	var shapeRef = map[reflect.Type]*ShapeInfo{}
	for _, opt := range opts {
		for _, shp := range opt.shapeUsage {
			if info, ok := shapeRef[shp.getType()]; ok {
				info.eNum += shp.base().executeNum
			} else {
				shapeInfo := &ShapeInfo{
					typ:    shp.getType(),
					eNum:   shp.base().executeNum,
					shapes: []IShapeGetter{shp},
				}
				shapeRef[shp.getType()] = shapeInfo
			}
		}
	}
	//sort
	o.shapeInfos = []*ShapeInfo{}
	for _, info := range shapeRef {
		o.shapeInfos = append(o.shapeInfos, info)
	}
	sort.Slice(o.shapeInfos, func(i, j int) bool {
		return o.shapeInfos[i].eNum > o.shapeInfos[j].eNum
	})

	o.lastCollectConsumption = time.Since(start)
}

func (o *optimizer) optimize(IdleTime time.Duration) {
	Log.Infof("start optimize, rest time: %v", IdleTime)
	o.expireTime = time.Now().Add(IdleTime)

	o.collect()

	rest := o.expire()
	Log.Infof("rest time: %v", rest)
}

func (o *optimizer) expire() time.Duration {
	return time.Until(o.expireTime)
}
