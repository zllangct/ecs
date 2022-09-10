package ecs

import (
	"reflect"
	"sort"
	"time"
)

type ShapeInfo struct {
	typ    reflect.Type
	eNum   int64
	shapes []IShape
}

type OptimizerReporter struct {
	shapeUsage map[reflect.Type]IShape
}

func (o *OptimizerReporter) init() {
	o.shapeUsage = map[reflect.Type]IShape{}
}

type optimizer struct {
	world                  *worldBase
	startTime              time.Time
	expireTime             time.Time
	lastSample             time.Time
	shapeInfos             []*ShapeInfo
	lastCollectConsumption time.Duration
}

func newOptimizer(world *worldBase) *optimizer {
	return &optimizer{world: world}
}

// Collect 采集分布在各个系统中的OptimizerReporter
func (o *optimizer) collect() {
	start := time.Now()
	var opts []*OptimizerReporter
	for _, value := range o.world.systemFlow.systems {
		system, ok := value.(ISystem)
		if !ok {
			continue
		}
		if system != nil {
			opts = append(opts, system.getOptimizer())
		}
	}
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
					shapes: []IShape{shp},
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

func (o *optimizer) optimize(IdleTime time.Duration, force bool) {
	Log.Infof("start optimize, rest time: %v", IdleTime)
	o.startTime = time.Now()
	o.lastSample = time.Now()
	o.expireTime = o.startTime.Add(IdleTime)

	o.collect()
	elapsed := o.elapsedStep()
	Log.Infof("collect step 1: %v", elapsed)

	o.memTidy(force)

	rest := o.expire()
	total := time.Now().Sub(o.startTime)
	Log.Infof("end optimize, rest time: %v, total: %v", rest, total)
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
	//seq := uint32(0)
	//m := map[interface{}][]*EntityInfo{}
	//o.world.entities.foreach(func(entity Entity, info *EntityInfo) bool {
	//	c := info.getCompound().Type()
	//	_, ok := m[c]
	//	if !ok {
	//		m[c] = []*EntityInfo{}
	//	}
	//	m[c] = append(m[c], info)
	//	return true
	//})
	//
	//elapsed := o.elapsedStep()
	//rest := o.expire()
	//Log.Infof("memTidy step 1: %v, expire: %v", elapsed, rest)
	//if !force && rest < time.Millisecond {
	//	return
	//}
	//
	//for _, infos := range m {
	//	for _, info := range infos {
	//		seq++
	//		for _, component := range info.components {
	//			component.setSeq(seq)
	//			c := o.world.components.getComponentSet(component.Type()).GetByEntity(int64(component.Owner().Entity()))
	//			verify := c.(IComponent)
	//			println(component.debugAddress(), verify.debugAddress())
	//			if verify.getSeq() != component.getSeq() {
	//				Log.Errorf("component seq error, %v, %v", verify.getSeq(), component.getSeq())
	//			}
	//		}
	//	}
	//}
	//
	//elapsed = o.elapsedStep()
	//rest = o.expire()
	//Log.Infof("memTidy step 2: %v, expire: %v", elapsed, rest)
	//if !force && rest < time.Millisecond {
	//	return
	//}
	//
	//for _, collection := range o.world.components.getCollections() {
	//	collection.Sort()
	//	if !force && o.expire() < time.Millisecond {
	//		break
	//	}
	//}
	//
	//elapsed = o.elapsedStep()
	//rest = o.expire()
	//Log.Infof("memTidy step 3: %v, expire: %v", elapsed, rest)
}
