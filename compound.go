package ecs

import (
	"reflect"
	"sort"
)

type Compound = OrderedIntSet[uint16]

func newCompoundFromMap(components map[reflect.Type]IComponent) Compound {
	c := make(Compound, len(components))
	i := 0
	for _, component := range components {
		c[i] = component.getIntType()
		i++
	}
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	return c
}
