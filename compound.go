package ecs

type Compound = OrderedIntSet[uint16]

func NewCompound(initCap ...int) Compound {
	cap := 0
	if len(initCap) > 0 {
		cap = initCap[0]
	}
	return make(Compound, 0, cap)
}
