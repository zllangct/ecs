package ecs

type EntityOption func(config *EntityConfig)

type EntityConfig struct {
	comps []Component
}

func (s *EntityConfig) initDefault() {
}

func WithComponents(comp ...Component) EntityOption {
	return func(config *EntityConfig) {
		config.comps = comp
	}
}

type EntityInfo struct {
	entity   Entity
	compound Compound
	world    *World
}

// Entity returns the entity ID
func (e *EntityInfo) Entity() Entity {
	return e.entity
}

func (e *EntityInfo) Add(comps ...Component) *EntityInfo {
	for _, comp := range comps {
		if e.compound.Exist(GetIntTypeByComp(comp)) {
			continue
		}
		if comp.IsNomadic() {
			continue
		}
		op := Operate{
			Entity: e.entity,
			Op:     ComponentOperateAdd,
			Comp:   comp,
		}
		e.world.componentOp(op)
	}
	return e
}

// Remove 移除组件（帧同步点统一生效）
func (e *EntityInfo) Remove(comps ...Component) *EntityInfo {
	for _, comp := range comps {
		if !e.compound.Exist(GetIntTypeByComp(comp)) {
			continue
		}
		op := Operate{
			Entity: e.entity,
			Op:     ComponentOperateDelete,
			Comp:   comp,
		}
		e.world.componentOp(op)
	}
	return e
}
