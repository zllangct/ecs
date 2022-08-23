package ecs

import "testing"

func TestEntityIDGenerator_NewID(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		e := NewEntityIDGenerator(10, 3)
		id1 := e.NewID()
		id2 := e.NewID()
		id3 := e.NewID()

		e.FreeID(id2)

		id4 := e.NewID()

		e.FreeID(id1)
		e.FreeID(id4)
		e.FreeID(id3)

		var m []int64
		for i := 0; i < 11; i++ {
			newID := e.NewID()
			m = append(m, newID)
		}

		for _, id := range m {
			e.FreeID(id)
		}
	})
}
