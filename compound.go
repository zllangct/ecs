package ecs

//go:generate go run ./cmd/ecs_internal_gen/main.go FixedCompound -p "ecs"

type FixedCompound interface {
	Compound() Compound
}

type Compound = OrderedIntSet[ComponentIntType]

func NewCompound(initCap ...int) Compound {
	capability := 0
	if len(initCap) > 0 {
		capability = initCap[0]
	}
	return make(Compound, 0, capability)
}
