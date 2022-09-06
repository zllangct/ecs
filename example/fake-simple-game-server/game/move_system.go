package game

import (
	"github.com/zllangct/ecs"
	"time"
)

type MoveSystemData struct {
	P *Position
	M *Movement
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	timeScale float64
	deltaTime time.Duration
	getter    *ecs.ShapeGetter[ecs.Shape2[Movement, Position], *ecs.Shape2[Movement, Position]]
}

func (m *MoveSystem) Init(initializer ecs.SystemInitializer) {
	m.SetRequirements(initializer, &Position{}, &Movement{})
	m.EventRegister("UpdateTimeScale", m.UpdateTimeScale)
	getter, err := ecs.NewShapeGetter[ecs.Shape2[Movement, Position]](m)
	if err != nil {
		ecs.Log.Error(err)
	}
	m.getter = getter
}

func (m *MoveSystem) UpdateTimeScale(timeScale []interface{}) {
	ecs.Log.Info("time scale change to ", timeScale[0])
	m.timeScale = timeScale[0].(float64)
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta

	m.deltaTime += delta
	isPrint := false
	if m.deltaTime > time.Second*3 {
		isPrint = true
		m.deltaTime = 0
	}

	iter := m.getter.Get()
	for shp := iter.Begin(); !iter.End(); shp = iter.Next() {
		mv := shp.C1
		p := shp.C2
		_, _ = p, mv
		p.X = p.X + int(float64(mv.Dir[0]*mv.V)*delta.Seconds())
		p.Y = p.Y + int(float64(mv.Dir[1]*mv.V)*delta.Seconds())
		p.Z = p.Z + int(float64(mv.Dir[2]*mv.V)*delta.Seconds())

		if isPrint {
			e := p.Owner().Entity()
			ecs.Log.Info("target id:", e, " delta:", delta, " current position:", p.X, p.Y, p.Z)
		}
	}
}
