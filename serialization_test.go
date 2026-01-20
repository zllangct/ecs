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

// Test component for world serialization tests
type testWorldComponent struct {
	Value int32
	Name  string
}

const (
	packetIdentifierTestWorldComponent  rockmem.PacketIdentifier = 65529
	packetIdentifierTestWorldComponent2 rockmem.PacketIdentifier = 65528
)

func (t *testWorldComponent) NewComponentSet() ComponentSet {
	return NewCSet[testWorldComponent]()
}

func (t *testWorldComponent) PacketIdentifier() rockmem.PacketIdentifier {
	return packetIdentifierTestWorldComponent
}

func (t *testWorldComponent) IsNomadic() bool {
	return false
}

func (t *testWorldComponent) IsDisposable() bool {
	return false
}

// Test component 2 for world serialization tests
type testWorldComponent2 struct {
	X float32
	Y float32
}

func (t *testWorldComponent2) NewComponentSet() ComponentSet {
	return NewCSet[testWorldComponent2]()
}

func (t *testWorldComponent2) PacketIdentifier() rockmem.PacketIdentifier {
	return packetIdentifierTestWorldComponent2
}

func (t *testWorldComponent2) IsNomadic() bool {
	return false
}

func (t *testWorldComponent2) IsDisposable() bool {
	return false
}

func init() {
	// Register test components for serialization
	RegisterComponent[testWorldComponent]((*testWorldComponent)(nil))
	RegisterComponent[testWorldComponent2]((*testWorldComponent2)(nil))
	RegisterComponent[dummyComponent]((*dummyComponent)(nil))
}

func TestEntityIDGeneratorSerialization(t *testing.T) {
	// Create an EntityIDGenerator with some allocated IDs
	original := NewEntityIDGenerator(100, 10)

	// Allocate some IDs
	ids := make([]Entity, 10)
	for i := 0; i < 10; i++ {
		ids[i] = original.NewID()
	}

	// Free some IDs to test the free list
	original.FreeID(ids[2])
	original.FreeID(ids[5])

	// Allocate more IDs
	for i := 0; i < 3; i++ {
		original.NewID()
	}

	// Serialize
	data := original.Marshal()

	// Validate serialized data
	if data.Len != original.len {
		t.Errorf("Len mismatch: expected %d, got %d", original.len, data.Len)
	}
	if data.Free != int32(original.free) {
		t.Errorf("Free mismatch: expected %d, got %d", original.free, data.Free)
	}
	if data.Pending != int32(original.pending) {
		t.Errorf("Pending mismatch: expected %d, got %d", original.pending, data.Pending)
	}

	// Deserialize
	restored := &EntityIDGenerator{}
	restored.Unmarshal(data)

	// Validate restored data
	if restored.len != original.len {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.len, restored.len)
	}
	if restored.free != original.free {
		t.Errorf("Restored Free mismatch: expected %d, got %d", original.free, restored.free)
	}
	if restored.pending != original.pending {
		t.Errorf("Restored Pending mismatch: expected %d, got %d", original.pending, restored.pending)
	}

	// Test that new allocations work correctly after restoration
	newID := restored.NewID()
	if newID == 0 {
		t.Error("Failed to allocate new ID after restoration")
	}
}

func TestEntityIDGeneratorSerializationWithRockmem(t *testing.T) {
	original := NewEntityIDGenerator(50, 5)

	// Allocate and free some IDs
	for i := 0; i < 20; i++ {
		original.NewID()
	}
	original.FreeID(Entity(5))
	original.FreeID(Entity(10))

	// Serialize to rockmem writer
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Create reader from written data
	reader := rockmem.NewReader(writer.Bytes())

	// Deserialize from rockmem reader
	restored := &EntityIDGenerator{}
	restored.UnmarshalFrom(reader)

	// Validate
	if restored.len != original.len {
		t.Errorf("Restored Len mismatch: expected %d, got %d", original.len, restored.len)
	}
}

func TestWorldSerialization(t *testing.T) {
	// Create a world with entities and components
	original := NewWorld().(*world)

	// Add some entities with components
	for i := 0; i < 5; i++ {
		entity := original.NewEntity()
		entity.Add(&testWorldComponent{Value: int32(i * 100), Name: "test"})
		if i%2 == 0 {
			entity.Add(&testWorldComponent2{X: float32(i), Y: float32(i * 2)})
		}
	}

	// Update frame counter
	original.frame = 42

	// Serialize
	data := original.Marshal()

	// Validate serialized data
	if data.Frame != 42 {
		t.Errorf("Frame mismatch: expected 42, got %d", data.Frame)
	}
	if len(data.Entities.EntityIds) != 5 {
		t.Errorf("Entity count mismatch: expected 5, got %d", len(data.Entities.EntityIds))
	}

	// Deserialize
	restored := &world{}
	restored.Unmarshal(data)

	// Validate restored data
	if restored.frame != original.frame {
		t.Errorf("Restored Frame mismatch: expected %d, got %d", original.frame, restored.frame)
	}
	if restored.entities.Len() != original.entities.Len() {
		t.Errorf("Restored entity count mismatch: expected %d, got %d",
			original.entities.Len(), restored.entities.Len())
	}

	// Validate component sets
	if len(restored.components) != len(original.components) {
		t.Errorf("Component set count mismatch: expected %d, got %d",
			len(original.components), len(restored.components))
	}

	// Check specific component data
	compType := GetIntType[testWorldComponent, *testWorldComponent]()
	origSet, origOk := original.components[compType]
	restoredSet, restoredOk := restored.components[compType]
	if origOk && restoredOk {
		if origSet.Len() != restoredSet.Len() {
			t.Errorf("testWorldComponent count mismatch: expected %d, got %d",
				origSet.Len(), restoredSet.Len())
		}
	} else if origOk != restoredOk {
		t.Errorf("testWorldComponent set existence mismatch: original=%v, restored=%v", origOk, restoredOk)
	}
}

func TestWorldSerializationWithRockmem(t *testing.T) {
	// Create a world with entities and components
	original := NewWorld().(*world)

	// Add some entities with components
	for i := 0; i < 10; i++ {
		entity := original.NewEntity()
		entity.Add(&testWorldComponent{Value: int32(i * 100), Name: "entity"})
		entity.Add(&testWorldComponent2{X: float32(i), Y: float32(i * 2)})
	}

	original.frame = 100

	// Serialize to rockmem writer
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	bytes := writer.Bytes()
	t.Logf("Serialized world with %d entities to %d bytes", original.entities.Len(), len(bytes))

	// Create reader from written data
	reader := rockmem.NewReader(bytes)

	// Deserialize from rockmem reader
	restored := &world{}
	restored.UnmarshalFrom(reader)

	// Validate restored data
	if restored.frame != original.frame {
		t.Errorf("Restored Frame mismatch: expected %d, got %d", original.frame, restored.frame)
	}
	if restored.entities.Len() != original.entities.Len() {
		t.Errorf("Restored entity count mismatch: expected %d, got %d",
			original.entities.Len(), restored.entities.Len())
	}
}

func TestWorldSerializationRoundTrip(t *testing.T) {
	// This test simulates a complete migration scenario:
	// 1. Create a world with state
	// 2. Serialize it
	// 3. Deserialize to a new world
	// 4. Verify all state is preserved

	// Create original world
	original := NewWorld().(*world)

	// Add entities with various components
	entities := make([]*EntityInfo, 20)
	for i := 0; i < 20; i++ {
		entities[i] = original.NewEntity()
		entities[i].Add(&testWorldComponent{Value: int32(i), Name: "test"})
		if i%3 == 0 {
			entities[i].Add(&testWorldComponent2{X: float32(i), Y: float32(i * 10)})
		}
	}

	// Simulate some frames
	original.frame = 500

	// Add disposable types for testing
	original.disposableTypes = []ComponentIntType{1, 2, 3}
	original.nomadicTypes = []ComponentIntType{4, 5}

	// Serialize to bytes (simulating network transfer)
	writer := rockmem.NewWriter()
	_, err := original.MarshalTo(writer)
	if err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}
	networkBytes := writer.Bytes()

	// "Transfer" over network and deserialize at destination
	reader := rockmem.NewReader(networkBytes)
	restored := &world{}
	restored.UnmarshalFrom(reader)

	// Verify all state
	if restored.frame != original.frame {
		t.Errorf("Frame not preserved: expected %d, got %d", original.frame, restored.frame)
	}

	if restored.entities.Len() != original.entities.Len() {
		t.Errorf("Entity count not preserved: expected %d, got %d",
			original.entities.Len(), restored.entities.Len())
	}

	// Verify disposable types
	if len(restored.disposableTypes) != len(original.disposableTypes) {
		t.Errorf("DisposableTypes count not preserved: expected %d, got %d",
			len(original.disposableTypes), len(restored.disposableTypes))
	}
	for i, dt := range original.disposableTypes {
		if i < len(restored.disposableTypes) && restored.disposableTypes[i] != dt {
			t.Errorf("DisposableType mismatch at %d: expected %d, got %d",
				i, dt, restored.disposableTypes[i])
		}
	}

	// Verify nomadic types
	if len(restored.nomadicTypes) != len(original.nomadicTypes) {
		t.Errorf("NomadicTypes count not preserved: expected %d, got %d",
			len(original.nomadicTypes), len(restored.nomadicTypes))
	}

	// Verify component data integrity
	compType := GetIntType[testWorldComponent, *testWorldComponent]()
	origSet, origExists := original.components[compType]
	restoredSet, restoredExists := restored.components[compType]

	if origExists != restoredExists {
		t.Errorf("Component set existence mismatch: original=%v, restored=%v", origExists, restoredExists)
	}

	if origExists && restoredExists {
		origCSet := origSet.(*CSet[testWorldComponent, *testWorldComponent])
		restoredCSet := restoredSet.(*CSet[testWorldComponent, *testWorldComponent])

		if origCSet.Len() != restoredCSet.Len() {
			t.Errorf("Component set size not preserved: expected %d, got %d",
				origCSet.Len(), restoredCSet.Len())
		}
	}
}

func TestEmptyWorldSerialization(t *testing.T) {
	original := NewWorld().(*world)

	// Serialize empty world
	data := original.Marshal()

	// Deserialize
	restored := &world{}
	restored.Unmarshal(data)

	// Verify
	if restored.entities.Len() != 0 {
		t.Errorf("Empty world serialization failed: expected 0 entities, got %d",
			restored.entities.Len())
	}
	if restored.frame != 0 {
		t.Errorf("Empty world frame should be 0, got %d", restored.frame)
	}
}

func TestNewWorldFromData(t *testing.T) {
	// Create and populate original world
	original := NewWorld().(*world)
	for i := 0; i < 5; i++ {
		entity := original.NewEntity()
		entity.Add(&testWorldComponent{Value: int32(i)})
	}
	original.frame = 123

	// Serialize
	data := original.Marshal()

	// Create new world from data
	restored := NewWorldFromData(data).(*world)

	// Verify
	if restored.frame != original.frame {
		t.Errorf("Frame mismatch: expected %d, got %d", original.frame, restored.frame)
	}
	if restored.entities.Len() != original.entities.Len() {
		t.Errorf("Entity count mismatch: expected %d, got %d",
			original.entities.Len(), restored.entities.Len())
	}
	if restored.status != WorldStatusInitialized {
		t.Errorf("Status should be Initialized, got %d", restored.status)
	}
}

func BenchmarkWorldSerialization(b *testing.B) {
	// Create a world with entities and components
	w := NewWorld().(*world)
	for i := 0; i < 1000; i++ {
		entity := w.NewEntity()
		entity.Add(&testWorldComponent{Value: int32(i), Name: "benchmark"})
		if i%2 == 0 {
			entity.Add(&testWorldComponent2{X: float32(i), Y: float32(i * 2)})
		}
	}
	w.frame = 1000

	b.Run("Marshal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = w.Marshal()
		}
	})

	data := w.Marshal()

	b.Run("Unmarshal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			restored := &world{}
			restored.Unmarshal(data)
		}
	})

	b.Run("MarshalTo", func(b *testing.B) {
		writer := rockmem.NewWriter()
		for n := 0; n < b.N; n++ {
			writer.Reset()
			_, _ = w.MarshalTo(writer)
		}
	})

	writer := rockmem.NewWriter()
	_, _ = w.MarshalTo(writer)
	bytes := writer.Bytes()

	b.Run("UnmarshalFrom", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			reader := rockmem.NewReader(bytes)
			restored := &world{}
			restored.UnmarshalFrom(reader)
		}
	})

	b.Logf("World with %d entities serialized to %d bytes", w.entities.Len(), len(bytes))
}
