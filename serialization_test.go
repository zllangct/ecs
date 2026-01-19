package ecs

import (
	"testing"

	rockmem "github.com/zllangct/rockmem/golang"
)

func TestUSetSerialization(t *testing.T) {
	// Create a USet with some test data
	original := NewUSet[dummyComponent](10)

	// Add some elements
	for i := 0; i < 5; i++ {
		item := &dummyComponent{Seq: int32(i * 10)}
		original.Add(item)
	}

	// Serialize
	data := original.Marshal()

	// Validate serialized data
	if data.EleSize != original.eleSize {
		t.Errorf("EleSize mismatch: expected %d, got %d", original.eleSize, data.EleSize)
	}
	if data.Len != original.len {
		t.Errorf("Len mismatch: expected %d, got %d", original.len, data.Len)
	}
	if data.InitSize != original.initSize {
		t.Errorf("InitSize mismatch: expected %d, got %d", original.initSize, data.InitSize)
	}

	// Deserialize to a new USet
	restored := NewUSet[dummyComponent]()
	restored.Unmarshal(data)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	for i := 0; i < original.Len(); i++ {
		origItem := original.Get(int64(i))
		restoredItem := restored.Get(int64(i))
		if origItem.Seq != restoredItem.Seq {
			t.Errorf("Data mismatch at index %d: expected %d, got %d", i, origItem.Seq, restoredItem.Seq)
		}
	}
}

func TestUSetSerializationWithRockmem(t *testing.T) {
	// Create a USet with some test data
	original := NewUSet[dummyComponent](10)

	// Add some elements
	for i := 0; i < 5; i++ {
		item := &dummyComponent{Seq: int32(i * 10)}
		original.Add(item)
	}

	// Serialize to rockmem writer
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Create reader from written data
	reader := rockmem.NewReader(writer.Bytes())

	// Deserialize from rockmem reader
	restored := NewUSet[dummyComponent]()
	restored.UnmarshalFrom(reader)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	for i := 0; i < original.Len(); i++ {
		origItem := original.Get(int64(i))
		restoredItem := restored.Get(int64(i))
		if origItem.Seq != restoredItem.Seq {
			t.Errorf("Data mismatch at index %d: expected %d, got %d", i, origItem.Seq, restoredItem.Seq)
		}
	}
}

func TestSparseArraySerialization(t *testing.T) {
	// Create a SparseArray with some test data
	original := NewSparseArray[EntityIndex, dummyComponent](10)

	// Add some elements at various keys
	keys := []EntityIndex{0, 5, 10, 25, 100}
	for i, key := range keys {
		item := &dummyComponent{Seq: int32(i * 10)}
		original.Add(key, item)
	}

	// Serialize
	data := original.Marshal()

	// Validate serialized data
	if data.MaxKey != int64(original.maxKey) {
		t.Errorf("MaxKey mismatch: expected %d, got %d", original.maxKey, data.MaxKey)
	}
	if data.ShrinkThreshold != original.shrinkThreshold {
		t.Errorf("ShrinkThreshold mismatch: expected %d, got %d", original.shrinkThreshold, data.ShrinkThreshold)
	}
	if data.IsKOrder != original.isKOrder {
		t.Errorf("IsKOrder mismatch: expected %v, got %v", original.isKOrder, data.IsKOrder)
	}

	// Deserialize to a new SparseArray
	restored := NewSparseArray[EntityIndex, dummyComponent]()
	restored.Unmarshal(data)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}
	if restored.maxKey != original.maxKey {
		t.Errorf("Restored MaxKey mismatch: expected %d, got %d", original.maxKey, restored.maxKey)
	}

	// Check all keys exist and have correct values
	for i, key := range keys {
		origItem := original.Get(key)
		restoredItem := restored.Get(key)
		if restoredItem == nil {
			t.Errorf("Key %d not found in restored SparseArray", key)
			continue
		}
		if origItem.Seq != restoredItem.Seq {
			t.Errorf("Data mismatch at key %d: expected %d, got %d", key, origItem.Seq, restoredItem.Seq)
		}
		_ = i
	}
}

func TestSparseArraySerializationWithRockmem(t *testing.T) {
	// Create a SparseArray with some test data
	original := NewSparseArray[EntityIndex, dummyComponent](10)

	// Add some elements at various keys
	keys := []EntityIndex{0, 5, 10, 25, 100}
	for i, key := range keys {
		item := &dummyComponent{Seq: int32(i * 10)}
		original.Add(key, item)
	}

	// Serialize to rockmem writer
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Create reader from written data
	reader := rockmem.NewReader(writer.Bytes())

	// Deserialize from rockmem reader
	restored := NewSparseArray[EntityIndex, dummyComponent]()
	restored.UnmarshalFrom(reader)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	// Check all keys exist and have correct values
	for i, key := range keys {
		origItem := original.Get(key)
		restoredItem := restored.Get(key)
		if restoredItem == nil {
			t.Errorf("Key %d not found in restored SparseArray", key)
			continue
		}
		if origItem.Seq != restoredItem.Seq {
			t.Errorf("Data mismatch at key %d: expected %d, got %d", key, origItem.Seq, restoredItem.Seq)
		}
		_ = i
	}
}

func TestCSetSerialization(t *testing.T) {
	// Create a CSet with some test data
	original := NewCSet[dummyComponent](10)

	// Add some components with entity IDs
	for i := 0; i < 5; i++ {
		entity := Entity(i * 10) // Entity IDs: 0, 10, 20, 30, 40
		comp := &dummyComponent{Seq: int32(i * 100)}
		original.Add(entity, comp)
	}

	// Serialize
	data := original.Marshal()

	// Deserialize to a new CSet
	restored := NewCSetFromData[dummyComponent](data)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	// Check all components
	for i := 0; i < 5; i++ {
		entity := Entity(i * 10)
		origComp := original.Get(entity)
		restoredComp := restored.Get(entity)
		if restoredComp == nil {
			t.Errorf("Entity %d component not found in restored CSet", entity)
			continue
		}
		origDummy := origComp.(*dummyComponent)
		restoredDummy := restoredComp.(*dummyComponent)
		if origDummy.Seq != restoredDummy.Seq {
			t.Errorf("Component mismatch for entity %d: expected %d, got %d",
				entity, origDummy.Seq, restoredDummy.Seq)
		}
	}
}

func TestCSetSerializationWithRockmem(t *testing.T) {
	// Create a CSet with some test data
	original := NewCSet[dummyComponent](10)

	// Add some components with entity IDs
	for i := 0; i < 5; i++ {
		entity := Entity(i * 10)
		comp := &dummyComponent{Seq: int32(i * 100)}
		original.Add(entity, comp)
	}

	// Serialize to rockmem writer
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Create reader from written data
	reader := rockmem.NewReader(writer.Bytes())

	// Deserialize from rockmem reader
	restored := NewCSet[dummyComponent]()
	restored.UnmarshalFrom(reader)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	// Check all components
	for i := 0; i < 5; i++ {
		entity := Entity(i * 10)
		origComp := original.Get(entity)
		restoredComp := restored.Get(entity)
		if restoredComp == nil {
			t.Errorf("Entity %d component not found in restored CSet", entity)
			continue
		}
		origDummy := origComp.(*dummyComponent)
		restoredDummy := restoredComp.(*dummyComponent)
		if origDummy.Seq != restoredDummy.Seq {
			t.Errorf("Component mismatch for entity %d: expected %d, got %d",
				entity, origDummy.Seq, restoredDummy.Seq)
		}
	}
}

func TestEmptyContainerSerialization(t *testing.T) {
	// Test empty USet
	t.Run("EmptyUSet", func(t *testing.T) {
		original := NewUSet[dummyComponent]()
		data := original.Marshal()

		restored := NewUSet[dummyComponent]()
		restored.Unmarshal(data)

		if restored.Len() != 0 {
			t.Errorf("Empty USet serialization failed: expected 0 length, got %d", restored.Len())
		}
	})

	// Test empty SparseArray
	t.Run("EmptySparseArray", func(t *testing.T) {
		original := NewSparseArray[EntityIndex, dummyComponent]()
		data := original.Marshal()

		restored := NewSparseArray[EntityIndex, dummyComponent]()
		restored.Unmarshal(data)

		if restored.Len() != 0 {
			t.Errorf("Empty SparseArray serialization failed: expected 0 length, got %d", restored.Len())
		}
	})

	// Test empty CSet
	t.Run("EmptyCSet", func(t *testing.T) {
		original := NewCSet[dummyComponent]()
		data := original.Marshal()

		restored := NewCSetFromData[dummyComponent](data)

		if restored.Len() != 0 {
			t.Errorf("Empty CSet serialization failed: expected 0 length, got %d", restored.Len())
		}
	})
}

func TestSerializationDataIntegrity(t *testing.T) {
	// Create a CSet with test data
	original := NewCSet[dummyComponent](100)

	// Add many components
	for i := 0; i < 100; i++ {
		entity := Entity(i)
		comp := &dummyComponent{Seq: int32(i)}
		original.Add(entity, comp)
	}

	// Serialize using rockmem
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	bytes := writer.Bytes()
	t.Logf("Serialized %d components to %d bytes", original.Len(), len(bytes))

	// Deserialize and verify all data
	reader := rockmem.NewReader(bytes)
	restored := NewCSet[dummyComponent]()
	restored.UnmarshalFrom(reader)

	if restored.Len() != original.Len() {
		t.Fatalf("Length mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	// Verify each component
	for i := 0; i < 100; i++ {
		entity := Entity(i)
		origComp := original.Get(entity)
		restoredComp := restored.Get(entity)

		if restoredComp == nil {
			t.Errorf("Entity %d not found in restored CSet", entity)
			continue
		}

		origDummy := origComp.(*dummyComponent)
		restoredDummy := restoredComp.(*dummyComponent)

		if origDummy.Seq != restoredDummy.Seq {
			t.Errorf("Data corruption at entity %d: expected %d, got %d",
				entity, origDummy.Seq, restoredDummy.Seq)
		}
	}
}

func BenchmarkCSetSerialization(b *testing.B) {
	// Create a CSet with test data
	c := NewCSet[dummyComponent](1000)
	for i := 0; i < 1000; i++ {
		entity := Entity(i)
		comp := &dummyComponent{Seq: int32(i)}
		c.Add(entity, comp)
	}

	b.Run("Marshal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = c.Marshal()
		}
	})

	data := c.Marshal()

	b.Run("Unmarshal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			restored := NewCSet[dummyComponent]()
			restored.Unmarshal(data)
		}
	})

	b.Run("MarshalTo", func(b *testing.B) {
		writer := rockmem.NewWriter()
		for n := 0; n < b.N; n++ {
			writer.Reset()
			_, _ = c.MarshalTo(writer)
		}
	})

	writer := rockmem.NewWriter()
	_, _ = c.MarshalTo(writer)
	bytes := writer.Bytes()

	b.Run("UnmarshalFrom", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			reader := rockmem.NewReader(bytes)
			restored := NewCSet[dummyComponent]()
			restored.UnmarshalFrom(reader)
		}
	})
}

func TestEntitySetSerialization(t *testing.T) {
	// Create an EntitySet with some test data
	original := NewEntitySet()

	// Add some entities
	for i := 0; i < 5; i++ {
		info := EntityInfo{
			entity:   Entity(int64(i) | (int64(i+1) << 32)), // Create entity with index and reuse
			compound: NewCompound(4),
			world:    nil,
		}
		// Add some component types
		info.compound = append(info.compound, ComponentIntType(i), ComponentIntType(i+10))
		original.Add(info)
	}

	// Serialize
	data := original.Marshal()

	// Validate serialized data
	if len(data.EntityIds) != original.Len() {
		t.Errorf("EntityIds count mismatch: expected %d, got %d", original.Len(), len(data.EntityIds))
	}

	// Deserialize
	restored := NewEntitySet()
	restored.Unmarshal(data)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	// Check each entity
	for i := 0; i < original.Len(); i++ {
		origInfo := original.SparseArray.Get(EntityIndex(i))
		restoredInfo := restored.SparseArray.Get(EntityIndex(i))

		if restoredInfo == nil {
			t.Errorf("Entity at index %d not found in restored EntitySet", i)
			continue
		}

		if origInfo.entity != restoredInfo.entity {
			t.Errorf("Entity mismatch at index %d: expected %d, got %d", i, origInfo.entity, restoredInfo.entity)
		}

		if len(origInfo.compound) != len(restoredInfo.compound) {
			t.Errorf("Compound length mismatch at index %d: expected %d, got %d",
				i, len(origInfo.compound), len(restoredInfo.compound))
			continue
		}

		for j, ct := range origInfo.compound {
			if restoredInfo.compound[j] != ct {
				t.Errorf("Compound type mismatch at entity %d, index %d: expected %d, got %d",
					i, j, ct, restoredInfo.compound[j])
			}
		}
	}
}

func TestEntitySetSerializationWithRockmem(t *testing.T) {
	// Create an EntitySet with some test data
	original := NewEntitySet()

	// Add some entities
	for i := 0; i < 10; i++ {
		info := EntityInfo{
			entity:   Entity(int64(i) | (int64(i+1) << 32)),
			compound: NewCompound(4),
			world:    nil,
		}
		info.compound = append(info.compound, ComponentIntType(i), ComponentIntType(i+10), ComponentIntType(i+20))
		original.Add(info)
	}

	// Serialize to rockmem writer
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	bytes := writer.Bytes()
	t.Logf("Serialized %d entities to %d bytes", original.Len(), len(bytes))

	// Create reader from written data
	reader := rockmem.NewReader(bytes)

	// Deserialize from rockmem reader
	restored := NewEntitySet()
	restored.UnmarshalFrom(reader)

	// Validate restored data
	if restored.Len() != original.Len() {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.Len(), restored.Len())
	}

	// Check each entity
	for i := 0; i < original.Len(); i++ {
		origInfo := original.SparseArray.Get(EntityIndex(i))
		restoredInfo := restored.SparseArray.Get(EntityIndex(i))

		if restoredInfo == nil {
			t.Errorf("Entity at index %d not found", i)
			continue
		}

		if origInfo.entity != restoredInfo.entity {
			t.Errorf("Entity mismatch at index %d: expected %d, got %d", i, origInfo.entity, restoredInfo.entity)
		}
	}
}

func TestEmptyEntitySetSerialization(t *testing.T) {
	original := NewEntitySet()
	data := original.Marshal()

	restored := NewEntitySetFromData(data)

	if restored.Len() != 0 {
		t.Errorf("Empty EntitySet serialization failed: expected 0 length, got %d", restored.Len())
	}
}

func BenchmarkEntitySetSerialization(b *testing.B) {
	// Create an EntitySet with test data
	es := NewEntitySet()
	for i := 0; i < 1000; i++ {
		info := EntityInfo{
			entity:   Entity(int64(i) | (int64(i+1) << 32)),
			compound: NewCompound(4),
			world:    nil,
		}
		info.compound = append(info.compound, ComponentIntType(i%100), ComponentIntType((i+10)%100))
		es.Add(info)
	}

	b.Run("Marshal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = es.Marshal()
		}
	})

	data := es.Marshal()

	b.Run("Unmarshal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			restored := NewEntitySet()
			restored.Unmarshal(data)
		}
	})

	b.Run("MarshalTo", func(b *testing.B) {
		writer := rockmem.NewWriter()
		for n := 0; n < b.N; n++ {
			writer.Reset()
			_, _ = es.MarshalTo(writer)
		}
	})

	writer := rockmem.NewWriter()
	_, _ = es.MarshalTo(writer)
	bytes := writer.Bytes()

	b.Run("UnmarshalFrom", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			reader := rockmem.NewReader(bytes)
			restored := NewEntitySet()
			restored.UnmarshalFrom(reader)
		}
	})
}
