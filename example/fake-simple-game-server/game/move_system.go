package game

import (
	"errors"
	"github.com/zllangct/ecs"
	"time"
)

type MoveSystemUtility struct {
	ecs.Utility[MoveSystemUtility]
}

func (m *MoveSystemUtility) UpdateTimeScale(scale float64) error {
	s := m.GetSystem()
	if s == nil {
		return errors.New("system is nil")
	}
	sys := s.(*MoveSystem)
	sys.timeScale = scale
	return nil
}

func (m *MoveSystemUtility) Move(entity ecs.Entity, v int, dir [3]int) error {
	s := m.GetSystem()
	if s == nil {
		return errors.New("system is nil")
	}
	mov := ecs.GetComponent[Movement](s, entity)
	mov.V = v
	mov.Dir = dir
	return nil
}

type MoveSystemData struct {
	P *Position
	M *Movement
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	timeScale float64
	deltaTime time.Duration
	getter    *ecs.Shape[MoveSystemData]
}

func (m *MoveSystem) Init(si ecs.SystemInitConstraint) {
	m.SetRequirements(si, &Position{}, &Movement{})
	ecs.BindUtility[MoveSystemUtility](si)
	m.getter = ecs.NewShape[MoveSystemData](si)
}

func (m *MoveSystem) UpdateTimeScale(timeScale []interface{}) error {
	ecs.Log.Info("time scale change to ", timeScale[0])
	m.timeScale = timeScale[0].(float64)
	return nil
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
		mv := shp.M
		p := shp.P
		_, _ = p, mv
		p.X = p.X + int(float64(mv.Dir[0]*mv.V)*delta.Seconds())
		p.Y = p.Y + int(float64(mv.Dir[1]*mv.V)*delta.Seconds())
		p.Z = p.Z + int(float64(mv.Dir[2]*mv.V)*delta.Seconds())

		if isPrint {
			e := p.Owner()
			ecs.Log.Info("target id:", e, " delta:", delta, " current position:", p.X, p.Y, p.Z)
		}
	}
}
