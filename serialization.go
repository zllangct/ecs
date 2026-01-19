package ecs

import (
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// USet serialization methods

// Marshal serializes USet[T] to SerializableUSetData
// It converts the typed data slice to raw bytes for storage
func (u *USet[T]) Marshal() *SerializableUSetData {
	data := &SerializableUSetData{
		EleSize:  u.eleSize,
		Len:      u.len,
		InitSize: u.initSize,
	}

	if u.len > 0 {
		// Convert typed data to raw bytes
		byteSize := int(u.len) * int(u.eleSize)
		data.Data = make([]byte, byteSize)
		srcPtr := unsafe.Pointer(&u.data[0])
		copy(data.Data, unsafe.Slice((*byte)(srcPtr), byteSize))
	}

	return data
}

// MarshalTo writes USet[T] to a rockmem.Writer
func (u *USet[T]) MarshalTo(writer rockmem.Writer) (uint, error) {
	data := u.Marshal()
	return data.WriteAsRoot(writer)
}

// Unmarshal deserializes SerializableUSetData to USet[T]
// IMPORTANT: The caller must ensure type T matches the original serialized type
func (u *USet[T]) Unmarshal(data *SerializableUSetData) {
	u.eleSize = data.EleSize
	u.len = data.Len
	u.initSize = data.InitSize

	if data.Len > 0 {
		// Calculate required capacity
		eleCount := int(data.Len)

		// Allocate new data slice
		u.data = u.a.alloc(eleCount, eleCount)

		// Copy raw bytes back to typed data
		byteSize := eleCount * int(u.eleSize)
		dstPtr := unsafe.Pointer(&u.data[0])
		copy(unsafe.Slice((*byte)(dstPtr), byteSize), data.Data[:byteSize])
	} else {
		u.data = u.a.alloc(0, int(u.initSize))
	}
}

// UnmarshalFrom reads USet[T] from a rockmem.Reader
func (u *USet[T]) UnmarshalFrom(reader *rockmem.Reader) {
	data := &SerializableUSetData{}
	data.ReadAsRoot(reader)
	u.Unmarshal(data)
}

// SparseArray serialization methods

// Marshal serializes SparseArray[K, V] to SerializableSparseArrayData
func (s *SparseArray[K, V]) Marshal() *SerializableSparseArrayData {
	data := &SerializableSparseArrayData{
		USetData:        *s.USet.Marshal(),
		MaxKey:          int64(s.maxKey),
		ShrinkThreshold: s.shrinkThreshold,
		InitSize:        int32(s.initSize),
		IsKOrder:        s.isKOrder,
	}

	// Copy indices slice
	if len(s.indices) > 0 {
		data.Indices = make([]int32, len(s.indices))
		copy(data.Indices, s.indices)
	}

	// Copy idx2Key slice
	if len(s.idx2Key) > 0 {
		data.Idx2Key = make([]int32, len(s.idx2Key))
		copy(data.Idx2Key, s.idx2Key)
	}

	return data
}

// MarshalTo writes SparseArray[K, V] to a rockmem.Writer
func (s *SparseArray[K, V]) MarshalTo(writer rockmem.Writer) (uint, error) {
	data := s.Marshal()
	return data.WriteAsRoot(writer)
}

// Unmarshal deserializes SerializableSparseArrayData to SparseArray[K, V]
// IMPORTANT: The caller must ensure types K and V match the original serialized types
func (s *SparseArray[K, V]) Unmarshal(data *SerializableSparseArrayData) {
	// Unmarshal the embedded USet
	s.USet.Unmarshal(&data.USetData)

	s.maxKey = K(data.MaxKey)
	s.shrinkThreshold = data.ShrinkThreshold
	s.initSize = int(data.InitSize)
	s.isKOrder = data.IsKOrder

	// Copy indices slice
	if len(data.Indices) > 0 {
		s.indices = make([]int32, len(data.Indices))
		copy(s.indices, data.Indices)
	} else {
		s.indices = make([]int32, 0, s.initSize)
	}

	// Copy idx2Key slice
	if len(data.Idx2Key) > 0 {
		s.idx2Key = make([]int32, len(data.Idx2Key))
		copy(s.idx2Key, data.Idx2Key)
	} else {
		s.idx2Key = []int32{}
	}
}

// UnmarshalFrom reads SparseArray[K, V] from a rockmem.Reader
func (s *SparseArray[K, V]) UnmarshalFrom(reader *rockmem.Reader) {
	data := &SerializableSparseArrayData{}
	data.ReadAsRoot(reader)
	s.Unmarshal(data)
}

// CSet serialization methods (inherits from SparseArray)

// Marshal serializes CSet[T] to SerializableSparseArrayData
// Since CSet embeds SparseArray, we just call the parent's Marshal method
func (c *CSet[T]) Marshal() *SerializableSparseArrayData {
	return c.SparseArray.Marshal()
}

// MarshalTo writes CSet[T] to a rockmem.Writer
func (c *CSet[T]) MarshalTo(writer rockmem.Writer) (uint, error) {
	return c.SparseArray.MarshalTo(writer)
}

// Unmarshal deserializes SerializableSparseArrayData to CSet[T]
// IMPORTANT: The caller must ensure type T matches the original serialized type
func (c *CSet[T]) Unmarshal(data *SerializableSparseArrayData) {
	c.SparseArray.Unmarshal(data)
}

// UnmarshalFrom reads CSet[T] from a rockmem.Reader
func (c *CSet[T]) UnmarshalFrom(reader *rockmem.Reader) {
	c.SparseArray.UnmarshalFrom(reader)
}

// NewCSetFromData creates a new CSet[T] from SerializableSparseArrayData
func NewCSetFromData[T ComponentObject](data *SerializableSparseArrayData) *CSet[T] {
	c := &CSet[T]{}
	c.Unmarshal(data)
	return c
}

// NewSparseArrayFromData creates a new SparseArray[K, V] from SerializableSparseArrayData
func NewSparseArrayFromData[K Integer, V any](data *SerializableSparseArrayData) *SparseArray[K, V] {
	s := &SparseArray[K, V]{}
	s.Unmarshal(data)
	return s
}

// NewUSetFromData creates a new USet[T] from SerializableUSetData
func NewUSetFromData[T any](data *SerializableUSetData) *USet[T] {
	u := &USet[T]{}
	u.Unmarshal(data)
	return u
}

// EntitySet serialization methods

// Marshal serializes EntitySet to SerializableEntitySetData
// Uses flattened parallel arrays for efficient serialization of variable-length compounds
func (es *EntitySet) Marshal() *SerializableEntitySetData {
	data := &SerializableEntitySetData{
		MaxKey:          int64(es.maxKey),
		ShrinkThreshold: es.shrinkThreshold,
		InitSize:        int32(es.initSize),
		IsKOrder:        es.isKOrder,
	}

	// Serialize all EntityInfos as parallel arrays
	count := es.Len()
	if count > 0 {
		data.EntityIds = make([]int64, count)
		data.CompoundOffsets = make([]int32, count)
		data.CompoundLengths = make([]int32, count)

		// First pass: calculate total compound data size
		totalCompoundSize := 0
		for i := 0; i < count; i++ {
			totalCompoundSize += len(es.data[i].compound)
		}

		data.CompoundData = make([]uint16, totalCompoundSize)

		// Second pass: fill the arrays
		compoundOffset := int32(0)
		for i := 0; i < count; i++ {
			info := &es.data[i]
			data.EntityIds[i] = info.entity.ToInt64()
			data.CompoundOffsets[i] = compoundOffset
			data.CompoundLengths[i] = int32(len(info.compound))

			// Copy compound data
			for j, ct := range info.compound {
				data.CompoundData[int(compoundOffset)+j] = uint16(ct)
			}
			compoundOffset += int32(len(info.compound))
		}
	}

	// Copy indices slice
	if len(es.indices) > 0 {
		data.Indices = make([]int32, len(es.indices))
		copy(data.Indices, es.indices)
	}

	// Copy idx2Key slice
	if len(es.idx2Key) > 0 {
		data.Idx2Key = make([]int32, len(es.idx2Key))
		copy(data.Idx2Key, es.idx2Key)
	}

	return data
}

// MarshalTo writes EntitySet to a rockmem.Writer
func (es *EntitySet) MarshalTo(writer rockmem.Writer) (uint, error) {
	data := es.Marshal()
	return data.WriteAsRoot(writer)
}

// Unmarshal deserializes SerializableEntitySetData to EntitySet
// Note: world pointer in each EntityInfo must be set separately after deserialization
func (es *EntitySet) Unmarshal(data *SerializableEntitySetData) {
	es.maxKey = EntityIndex(data.MaxKey)
	es.shrinkThreshold = data.ShrinkThreshold
	es.initSize = int(data.InitSize)
	es.isKOrder = data.IsKOrder

	// Deserialize all EntityInfos from parallel arrays
	count := len(data.EntityIds)
	if count > 0 {
		// Initialize USet fields manually since we're not using raw bytes
		es.USet.eleSize = uint64(unsafe.Sizeof(EntityInfo{}))
		es.USet.len = int64(count)
		es.USet.initSize = int64(es.initSize)
		es.USet.data = es.USet.a.alloc(count, count)

		for i := 0; i < count; i++ {
			es.data[i].entity = Entity(data.EntityIds[i])
			es.data[i].world = nil // Must be set by caller

			// Reconstruct compound from flattened data
			offset := data.CompoundOffsets[i]
			length := data.CompoundLengths[i]
			if length > 0 {
				es.data[i].compound = make(Compound, length)
				for j := int32(0); j < length; j++ {
					es.data[i].compound[j] = ComponentIntType(data.CompoundData[offset+j])
				}
			} else {
				es.data[i].compound = NewCompound(4)
			}
		}
	} else {
		es.USet = *NewUSet[EntityInfo](es.initSize)
	}

	// Copy indices slice
	if len(data.Indices) > 0 {
		es.indices = make([]int32, len(data.Indices))
		copy(es.indices, data.Indices)
	} else {
		es.indices = make([]int32, 0, es.initSize)
	}

	// Copy idx2Key slice
	if len(data.Idx2Key) > 0 {
		es.idx2Key = make([]int32, len(data.Idx2Key))
		copy(es.idx2Key, data.Idx2Key)
	} else {
		es.idx2Key = []int32{}
	}
}

// UnmarshalFrom reads EntitySet from a rockmem.Reader
func (es *EntitySet) UnmarshalFrom(reader *rockmem.Reader) {
	data := &SerializableEntitySetData{}
	data.ReadAsRoot(reader)
	es.Unmarshal(data)
}

// SetWorldForAll sets the world pointer for all EntityInfos after deserialization
func (es *EntitySet) SetWorldForAll(w *world) {
	for i := 0; i < es.Len(); i++ {
		es.data[i].SetWorld(w)
	}
}

// NewEntitySetFromData creates a new EntitySet from SerializableEntitySetData
func NewEntitySetFromData(data *SerializableEntitySetData) *EntitySet {
	es := &EntitySet{}
	es.Unmarshal(data)
	return es
}

// EntityInfo helper methods for setting world after deserialization

// SetWorld sets the world pointer after deserialization
func (e *EntityInfo) SetWorld(w *world) {
	e.world = w
}
