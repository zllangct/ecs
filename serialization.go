package ecs

import (
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"sort"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// SnapshotVersion 快照格式版本（R19）：不兼容演进递增，恢复时校验拒绝
const SnapshotVersion uint64 = 1

type validatableRoot interface {
	ValidateAsRoot(*rockmem.Reader) bool
}

// validateSparseMetadata 校验快照稀疏表元数据（Len、Indices、Idx2Key 互逆一致性与 IsKOrder 真实性），
// 并返回用于恢复 lastKey 的最大/末位键，防止读档后再插入触发越界或 merge-join 漏行。
func validateSparseMetadata(n int, indices, idx2Key []int32, isKOrder bool) (int32, error) {
	if len(idx2Key) != n {
		return 0, fmt.Errorf("ecs: sparse idx2Key length mismatch: want %d, got %d", n, len(idx2Key))
	}
	for i, k := range idx2Key {
		if k < 0 || int(k) >= len(indices) {
			return 0, fmt.Errorf("ecs: sparse idx2Key[%d]=%d out of indices range [0,%d)", i, k, len(indices))
		}
		if indices[k] != int32(i+1) {
			return 0, fmt.Errorf("ecs: sparse reciprocal index mismatch at key %d: got %d, want %d", k, indices[k], i+1)
		}
		if isKOrder && i > 0 && k <= idx2Key[i-1] {
			return 0, fmt.Errorf("ecs: sparse IsKOrder=true but idx2Key not strictly ascending at %d (%d <= %d)", i, k, idx2Key[i-1])
		}
	}
	nonZero := 0
	for _, v := range indices {
		if v < 0 || int(v) > n {
			return 0, fmt.Errorf("ecs: sparse indices value %d out of range [0,%d]", v, n)
		}
		if v != 0 {
			nonZero++
		}
	}
	if nonZero != n {
		return 0, fmt.Errorf("ecs: sparse indices active count %d != len %d", nonZero, n)
	}
	if n == 0 || !isKOrder {
		return 0, nil
	}
	return idx2Key[n-1], nil
}

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

// Unmarshal deserializes SerializableUSetData to USet[T].
// R19：无论 Len 是否为 0，均校验快照与当前类型的布局一致性（EleSize == sizeof(T)）、非负长度与 InitSize，
// 非法输入返回 error 且不修改目标（原子失败）。
func (u *USet[T]) Unmarshal(data *SerializableUSetData) error {
	wantSize := uint64(TypeOf[T]().Size())
	if data.Len < 0 {
		return fmt.Errorf("ecs: USet negative length %d", data.Len)
	}
	if data.InitSize < 0 {
		return fmt.Errorf("ecs: USet negative initSize %d", data.InitSize)
	}
	if data.EleSize != wantSize {
		return fmt.Errorf("ecs: USet element size mismatch: snapshot=%d current=%d", data.EleSize, wantSize)
	}
	if data.Len > 0 {
		eleCount := int(data.Len)
		if int64(eleCount) != data.Len {
			return fmt.Errorf("ecs: USet invalid length %d", data.Len)
		}
		if wantSize > 0 && uint64(eleCount) > uint64(math.MaxInt)/wantSize {
			return fmt.Errorf("ecs: USet byte size overflow for len %d", data.Len)
		}
		byteSize := eleCount * int(wantSize)
		if byteSize > len(data.Data) {
			return fmt.Errorf("ecs: USet data too short: need %d bytes, got %d", byteSize, len(data.Data))
		}

		u.data = u.a.alloc(eleCount, eleCount)
		if byteSize > 0 {
			dstPtr := unsafe.Pointer(&u.data[0])
			copy(unsafe.Slice((*byte)(dstPtr), byteSize), data.Data[:byteSize])
		}
	} else {
		u.data = u.a.alloc(0, int(data.InitSize))
	}
	u.eleSize = wantSize
	u.len = data.Len
	u.initSize = data.InitSize
	return nil
}

// UnmarshalFrom reads USet[T] from a rockmem.Reader
func (u *USet[T]) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableUSetData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableUSetData root buffer")
	}
	data.ReadAsRoot(reader)
	return u.Unmarshal(data)
}

// SparseArray serialization methods

// Marshal serializes SparseArray[K, V] to SerializableSparseArrayData
func (s *SparseArray[K, V]) Marshal() *SerializableSparseArrayData {
	return s.marshalWithData(*s.USet.Marshal())
}

func (s *SparseArray[K, V]) marshalWithData(elements SerializableUSetData) *SerializableSparseArrayData {
	data := &SerializableSparseArrayData{
		USetData:        elements,
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
func (s *SparseArray[K, V]) Unmarshal(data *SerializableSparseArrayData) error {
	if data.MaxKey < 0 {
		return fmt.Errorf("ecs: SparseArray negative maxKey %d", data.MaxKey)
	}
	lastKey, err := validateSparseMetadata(int(data.USetData.Len), data.Indices, data.Idx2Key, data.IsKOrder)
	if err != nil {
		return err
	}
	// Unmarshal the embedded USet
	if err := s.USet.Unmarshal(&data.USetData); err != nil {
		return err
	}

	s.maxKey = K(data.MaxKey)
	s.shrinkThreshold = data.ShrinkThreshold
	s.initSize = int(data.InitSize)
	s.isKOrder = data.IsKOrder
	s.lastKey = K(lastKey)

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
	return nil
}

// UnmarshalFrom reads SparseArray[K, V] from a rockmem.Reader
func (s *SparseArray[K, V]) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableSparseArrayData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableSparseArrayData root buffer")
	}
	data.ReadAsRoot(reader)
	return s.Unmarshal(data)
}

// CSet serialization methods (inherits from SparseArray)

// rockmemMessage 是 rockmem 生成组件的序列化接口
type rockmemMessage interface {
	WriteAsRoot(rockmem.Writer) (uint, error)
	ReadAsRoot(*rockmem.Reader)
}

// Marshal serializes CSet[T] to SerializableSparseArrayData.
// 组件实现 rockmem.Message 时走逐元素 IDL 序列化（Data 布局：
// [(Len+1) int64 偏移表][逐元素 blob 拼接]，EleSize=0 作标记），
// 支持 string/slice/嵌套等非 POD 字段；否则回退 legacy raw 布局（仅 POD 安全）。
// R19：IDL 序列化失败返回 error，不再静默回退 raw（raw 含指针的快照跨进程无效）
func (c *CSet[T, TP]) Marshal() (*SerializableSparseArrayData, error) {
	var probe TP
	if _, ok := any(probe).(rockmemMessage); !ok || c.Len() == 0 {
		return c.SparseArray.Marshal(), nil
	}

	n := c.Len()
	if n > math.MaxInt/8-1 {
		return nil, fmt.Errorf("ecs: IDL offset table too large for %d components", n)
	}
	headerSize := (n + 1) * 8
	var encoded []byte
	w := rockmem.NewWriter()
	for i := 0; i < n; i++ {
		w.Reset()
		m := any(TP(&c.data[i])).(rockmemMessage)
		if _, err := m.WriteAsRoot(w); err != nil {
			return nil, fmt.Errorf("ecs: marshal component %d: %w", i, err)
		}
		elem := w.Bytes()
		if i == 0 {
			// 对固定宽度组件一次分配。限制变长编码的估计，避免超大的首元素放大整池容量。
			width := min(len(elem), int(c.eleSize))
			capacity := headerSize
			if width <= (math.MaxInt-headerSize)/n {
				capacity += n * width
			}
			encoded = make([]byte, headerSize, capacity)
		}
		encoded = append(encoded, elem...)
		binary.LittleEndian.PutUint64(encoded[(i+1)*8:], uint64(len(encoded)-headerSize))
	}
	return c.SparseArray.marshalWithData(SerializableUSetData{
		EleSize:  0, // IDL 路径标记；偏移相对于首个 blob
		Len:      c.len,
		InitSize: c.USet.initSize,
		Data:     encoded,
	}), nil
}

// MarshalTo writes CSet[T] to a rockmem.Writer
func (c *CSet[T, TP]) MarshalTo(writer rockmem.Writer) (uint, error) {
	data, err := c.Marshal()
	if err != nil {
		return 0, err
	}
	return data.WriteAsRoot(writer)
}

// Unmarshal deserializes SerializableSparseArrayData to CSet[T]
// IMPORTANT: The caller must ensure type T matches the original serialized type
func (c *CSet[T, TP]) Unmarshal(data *SerializableSparseArrayData) error {
	if data.USetData.EleSize == 0 && data.USetData.Len > 0 {
		return c.unmarshalIDL(data)
	}
	return c.SparseArray.Unmarshal(data)
}

// unmarshalIDL 按 IDL 布局（偏移表 + 逐元素 blob）恢复组件。
// R19：校验偏移表非负/单调/在数据界内、稀疏元数据一致，以及元素 blob 的 ValidateAsRoot，损坏快照原子失败。
func (c *CSet[T, TP]) unmarshalIDL(data *SerializableSparseArrayData) error {
	if data.USetData.Len <= 0 || data.USetData.InitSize < 0 || data.MaxKey < 0 {
		return fmt.Errorf("ecs: invalid IDL USetData header (len=%d init=%d maxKey=%d)",
			data.USetData.Len, data.USetData.InitSize, data.MaxKey)
	}
	n := int(data.USetData.Len)
	if int64(n) != data.USetData.Len || n > (math.MaxInt/8)-1 {
		return fmt.Errorf("ecs: IDL length overflow %d", data.USetData.Len)
	}
	lastKey, err := validateSparseMetadata(n, data.Indices, data.Idx2Key, data.IsKOrder)
	if err != nil {
		return err
	}

	raw := data.USetData.Data
	hdrLen := int64(n+1) * 8
	if int64(len(raw)) < hdrLen {
		return fmt.Errorf("ecs: IDL snapshot header truncated: need %d, got %d", hdrLen, len(raw))
	}
	var lastOffset uint64
	for i := 0; i <= n; i++ {
		off := binary.LittleEndian.Uint64(raw[i*8:])
		if off > uint64(math.MaxInt64) {
			return fmt.Errorf("ecs: IDL offset overflow at %d", i)
		}
		if off < lastOffset {
			return fmt.Errorf("ecs: IDL offset table not monotonic at %d", i)
		}
		lastOffset = off
	}
	blobLen := int64(len(raw)) - hdrLen
	if lastOffset > uint64(blobLen) {
		return fmt.Errorf("ecs: IDL blob truncated: need %d, got %d", lastOffset, blobLen)
	}

	// 先对临时缓冲区解码与 ValidateAsRoot，避免破坏目标 c（原子失败）
	tmpData := c.USet.a.alloc(n, n)
	// Custom decoders may retain their reader. Give each row a distinct
	// reader, and bound the retained block when just one row survives.
	var readers []rockmem.Reader
	blobStart := int(hdrLen) + int(binary.LittleEndian.Uint64(raw))
	for i := 0; i < n; i++ {
		blobEnd := int(hdrLen) + int(binary.LittleEndian.Uint64(raw[(i+1)*8:]))
		blob := raw[blobStart:blobEnd]
		blobStart = blobEnd
		m, ok := any(TP(&tmpData[i])).(rockmemMessage)
		if !ok {
			return fmt.Errorf("ecs: component type %T does not implement rockmemMessage for IDL restore", TP(nil))
		}
		if len(readers) == 0 {
			readers = make([]rockmem.Reader, min(128, n-i))
		}
		r := &readers[0]
		readers = readers[1:]
		*r = *rockmem.NewReader(blob)
		if vr, ok := any(m).(validatableRoot); ok {
			if !vr.ValidateAsRoot(r) {
				return fmt.Errorf("ecs: IDL component %d failed ValidateAsRoot", i)
			}
		}
		m.ReadAsRoot(r)
	}

	c.USet.eleSize = uint64(TypeOf[T]().Size())
	c.USet.len = int64(n)
	c.USet.initSize = data.USetData.InitSize
	c.USet.data = tmpData

	// SparseArray 元数据（与 SparseArray.Unmarshal 相同部分）
	c.maxKey = EntityIndex(data.MaxKey)
	c.shrinkThreshold = data.ShrinkThreshold
	c.initSize = int(data.InitSize)
	c.isKOrder = data.IsKOrder
	c.lastKey = EntityIndex(lastKey)
	if len(data.Indices) > 0 {
		c.indices = make([]int32, len(data.Indices))
		copy(c.indices, data.Indices)
	} else {
		c.indices = make([]int32, 0, c.initSize)
	}
	if len(data.Idx2Key) > 0 {
		c.idx2Key = make([]int32, len(data.Idx2Key))
		copy(c.idx2Key, data.Idx2Key)
	} else {
		c.idx2Key = []int32{}
	}
	return nil
}

// UnmarshalFrom reads CSet[T] from a rockmem.Reader
func (c *CSet[T, TP]) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableSparseArrayData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableSparseArrayData root buffer")
	}
	data.ReadAsRoot(reader)
	return c.Unmarshal(data)
}

// NewCSetFromData creates a new CSet[T] from SerializableSparseArrayData
func NewCSetFromData[T any, TP ComponentPointer[T]](data *SerializableSparseArrayData) (*CSet[T, TP], error) {
	c := &CSet[T, TP]{}
	if err := c.Unmarshal(data); err != nil {
		return nil, err
	}
	return c, nil
}

// NewSparseArrayFromData creates a new SparseArray[K, V] from SerializableSparseArrayData
func NewSparseArrayFromData[K Integer, V any](data *SerializableSparseArrayData) (*SparseArray[K, V], error) {
	s := &SparseArray[K, V]{}
	if err := s.Unmarshal(data); err != nil {
		return nil, err
	}
	return s, nil
}

// NewUSetFromData creates a new USet[T] from SerializableUSetData
func NewUSetFromData[T any](data *SerializableUSetData) (*USet[T], error) {
	u := &USet[T]{}
	if err := u.Unmarshal(data); err != nil {
		return nil, err
	}
	return u, nil
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

		data.CompoundData = make([]uint64, totalCompoundSize)

		// Second pass: fill the arrays
		compoundOffset := int32(0)
		for i := 0; i < count; i++ {
			info := &es.data[i]
			data.EntityIds[i] = info.entity.ToInt64()
			data.CompoundOffsets[i] = compoundOffset
			data.CompoundLengths[i] = int32(len(info.compound))

			// Copy compound data
			for j, ct := range info.compound {
				data.CompoundData[int(compoundOffset)+j] = uint64(ct)
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
// R19：校验 parallel arrays 一致性（offsets/lengths 与 CompoundData 界内），非法返回 error
func (es *EntitySet) Unmarshal(data *SerializableEntitySetData) error {
	if data.InitSize < 0 || data.MaxKey < 0 {
		return fmt.Errorf("ecs: entity set negative initSize=%d or maxKey=%d", data.InitSize, data.MaxKey)
	}
	count := len(data.EntityIds)
	lastKey, err := validateSparseMetadata(count, data.Indices, data.Idx2Key, data.IsKOrder)
	if err != nil {
		return err
	}
	if count > 0 {
		if len(data.CompoundOffsets) != count || len(data.CompoundLengths) != count {
			return fmt.Errorf("ecs: entity parallel arrays length mismatch: ids=%d offsets=%d lengths=%d",
				count, len(data.CompoundOffsets), len(data.CompoundLengths))
		}
		prevEnd := int32(0)
		for i := 0; i < count; i++ {
			if Entity(data.EntityIds[i]).Index() != EntityIndex(data.Idx2Key[i]) {
				return fmt.Errorf("ecs: entity[%d] index %d != idx2Key %d", i, Entity(data.EntityIds[i]).Index(), data.Idx2Key[i])
			}
			off, ln := data.CompoundOffsets[i], data.CompoundLengths[i]
			if off < 0 || ln < 0 || off != prevEnd || int64(off)+int64(ln) > int64(len(data.CompoundData)) {
				return fmt.Errorf("ecs: entity compound span invalid at %d (off=%d len=%d data=%d)", i, off, ln, len(data.CompoundData))
			}
			prevEnd = off + ln
		}
	}

	es.maxKey = EntityIndex(data.MaxKey)
	es.shrinkThreshold = data.ShrinkThreshold
	es.initSize = int(data.InitSize)
	es.isKOrder = data.IsKOrder
	es.lastKey = EntityIndex(lastKey)

	// Deserialize all EntityInfos from parallel arrays
	if count > 0 {
		// Initialize USet fields manually since we're not using raw bytes
		es.USet.eleSize = uint64(unsafe.Sizeof(EntityInfo{}))
		es.USet.len = int64(count)
		es.USet.initSize = int64(es.initSize)
		es.USet.data = es.USet.a.alloc(count, count)

		for first := 0; first < count; {
			stop := first + min(128, count-first)
			base := int(data.CompoundOffsets[first])
			end := int(data.CompoundOffsets[stop-1]) + int(data.CompoundLengths[stop-1])
			block := make(Compound, end-base)
			for i, typ := range data.CompoundData[base:end] {
				block[i] = ComponentIntType(typ)
			}
			for i := first; i < stop; i++ {
				lo := int(data.CompoundOffsets[i]) - base
				hi := lo + int(data.CompoundLengths[i])
				info := &es.data[i]
				info.entity = Entity(data.EntityIds[i])
				info.world = nil // Must be set by caller.
				// Appending to one entity must never overwrite its neighbor.
				info.compound = block[lo:hi:hi]
				if lo == hi {
					info.compound = NewCompound(4)
				}
			}
			first = stop
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
	return nil
}

// UnmarshalFrom reads EntitySet from a rockmem.Reader
func (es *EntitySet) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableEntitySetData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableEntitySetData root buffer")
	}
	data.ReadAsRoot(reader)
	return es.Unmarshal(data)
}

// SetWorldForAll sets the world pointer for all EntityInfos after deserialization
func (es *EntitySet) SetWorldForAll(w *World) {
	for i := 0; i < es.Len(); i++ {
		es.data[i].SetWorld(w)
	}
}

// NewEntitySetFromData creates a new EntitySet from SerializableEntitySetData
func NewEntitySetFromData(data *SerializableEntitySetData) (*EntitySet, error) {
	es := &EntitySet{}
	if err := es.Unmarshal(data); err != nil {
		return nil, err
	}
	return es, nil
}

// EntityInfo helper methods for setting world after deserialization

// SetWorld sets the world pointer after deserialization
func (e *EntityInfo) SetWorld(w *World) {
	e.world = w
}

// EntityIDGenerator serialization methods

// Marshal serializes EntityIDGenerator to SerializableEntityIDGeneratorData
func (e *EntityIDGenerator) Marshal() *SerializableEntityIDGeneratorData {
	data := &SerializableEntityIDGeneratorData{
		Free:      int32(e.free),
		Pending:   int32(e.pending),
		Len:       e.len,
		DelayFree: e.delayFree,
		DelayCap:  e.delayCap,
	}

	// Convert ReuseID slice to int64 slice
	if len(e.ids) > 0 {
		data.Ids = make([]int64, len(e.ids))
		for i, id := range e.ids {
			data.Ids[i] = id.ToInt64()
		}
	}

	// Convert removeDelay slice to int64 slice
	if len(e.removeDelay) > 0 {
		data.RemoveDelay = make([]int64, len(e.removeDelay))
		for i, id := range e.removeDelay {
			data.RemoveDelay[i] = id.ToInt64()
		}
	}

	return data
}

// MarshalTo writes EntityIDGenerator to a rockmem.Writer
func (e *EntityIDGenerator) MarshalTo(writer rockmem.Writer) (uint, error) {
	data := e.Marshal()
	return data.WriteAsRoot(writer)
}

// Unmarshal deserializes SerializableEntityIDGeneratorData to EntityIDGenerator.
// R19：校验 freelist 不变量（pending/free/len/delayFree 不越界），非法返回 error
func (e *EntityIDGenerator) Unmarshal(data *SerializableEntityIDGeneratorData) error {
	if data.Pending < 0 || int64(len(data.Ids)) < int64(data.Pending) {
		return fmt.Errorf("ecs: id generator pending %d exceeds ids %d", data.Pending, len(data.Ids))
	}
	if data.Len < 0 || data.DelayFree < 0 || data.DelayCap < 0 {
		return fmt.Errorf("ecs: id generator negative field (len=%d delayFree=%d delayCap=%d)",
			data.Len, data.DelayFree, data.DelayCap)
	}
	delayCap := data.DelayCap
	if delayCap <= 0 {
		delayCap = 10 // 与 NewEntityIDGenerator 默认值对齐
	}
	if data.DelayFree > delayCap {
		return fmt.Errorf("ecs: id generator delayFree %d exceeds delayCap %d", data.DelayFree, delayCap)
	}
	if len(data.Ids) > 0 {
		if data.Free < 0 || (data.Free != data.Pending && int(data.Free) >= len(data.Ids)) {
			return fmt.Errorf("ecs: id generator free %d out of range [0,%d)", data.Free, len(data.Ids))
		}
	} else if data.Len != 0 || data.DelayFree != 0 {
		return fmt.Errorf("ecs: id generator empty ids with non-zero len=%d or delayFree=%d", data.Len, data.DelayFree)
	}

	// Derived merge hints refer to slots in the old freelist.
	e.mergeHint = 0
	e.freePages = nil
	e.free = EntityIndex(data.Free)
	e.pending = EntityIndex(data.Pending)
	e.len = data.Len
	e.delayFree = data.DelayFree
	e.delayCap = delayCap

	// Convert int64 slice to ReuseID slice
	if len(data.Ids) > 0 {
		e.ids = make([]ReuseID, len(data.Ids))
		for i, id := range data.Ids {
			e.ids[i] = Entity(id).toReuseID()
		}
	} else {
		// 空数据等价于全新生成器：维持不变量 pending==len(ids)、index 0 保留不用、
		// ids[i].index=i+1 构成初始 freelist 链
		e.ids = make([]ReuseID, 1)
		e.ids[0].index = 1
		e.free = 1
		e.pending = 1
	}

	// Convert int64 slice to removeDelay slice
	if len(data.RemoveDelay) > 0 {
		if int(e.delayFree) > len(data.RemoveDelay) {
			return fmt.Errorf("ecs: id generator delayFree %d exceeds removeDelay len %d", e.delayFree, len(data.RemoveDelay))
		}
		e.removeDelay = make([]ReuseID, len(data.RemoveDelay))
		for i, id := range data.RemoveDelay {
			e.removeDelay[i] = Entity(id).toReuseID()
		}
		if int32(len(e.removeDelay)) < e.delayCap {
			e.removeDelay = append(e.removeDelay, make([]ReuseID, int(e.delayCap)-len(e.removeDelay))...)
		}
	} else {
		if e.delayFree > 0 {
			return fmt.Errorf("ecs: id generator delayFree %d with empty removeDelay", e.delayFree)
		}
		e.removeDelay = make([]ReuseID, e.delayCap)
	}
	return nil
}

// UnmarshalFrom reads EntityIDGenerator from a rockmem.Reader
func (e *EntityIDGenerator) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableEntityIDGeneratorData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableEntityIDGeneratorData root buffer")
	}
	data.ReadAsRoot(reader)
	return e.Unmarshal(data)
}

// NewEntityIDGeneratorFromData creates a new EntityIDGenerator from SerializableEntityIDGeneratorData
func NewEntityIDGeneratorFromData(data *SerializableEntityIDGeneratorData) (*EntityIDGenerator, error) {
	e := &EntityIDGenerator{}
	if err := e.Unmarshal(data); err != nil {
		return nil, err
	}
	return e, nil
}

// World serialization methods

// Marshal serializes serializableWorld to SerializableWorldData
func (w *serializableWorld) Marshal() (*SerializableWorldData, error) {
	return w.marshalWithArchetypes(nil)
}

// marshalWithArchetypes 只读联合导出世界状态（含可选组表列与 CSet 残段），按 ComponentIntType 升序确定性输出。
func (w *serializableWorld) marshalWithArchetypes(archetypes *ArchetypeRegistry) (*SerializableWorldData, error) {
	data := &SerializableWorldData{
		Version: SnapshotVersion,
		Frame:   w.frame,
	}

	// Serialize entities
	if w.entities != nil {
		data.Entities = *w.entities.Marshal()
	}

	// Serialize ID generator
	if w.idGenerator != nil {
		data.IdGenerator = *w.idGenerator.Marshal()
	}

	// 汇总待导出组件类型并按升序确定性导出
	typeSet := make(map[ComponentIntType]struct{}, len(w.components))
	for compType := range w.components {
		typeSet[compType] = struct{}{}
	}
	if archetypes != nil {
		for compType := range archetypes.ownerOf {
			typeSet[compType] = struct{}{}
		}
	}
	if len(typeSet) > 0 {
		compTypes := make([]ComponentIntType, 0, len(typeSet))
		for compType := range typeSet {
			compTypes = append(compTypes, compType)
		}
		sort.Slice(compTypes, func(i, j int) bool { return compTypes[i] < compTypes[j] })

		data.ComponentSets = make([]SerializableComponentSetEntry, 0, len(compTypes))
		registry := GetComponentRegistry()
		for _, compType := range compTypes {
			compSet := w.components[compType]
			var a *Archetype
			var col int
			if archetypes != nil {
				if owner, ok := archetypes.owner(compType); ok {
					a = owner
					col, _ = a.colOf(compType)
				}
			}
			if compSet == nil && a != nil && a.Len() > 0 {
				var ok bool
				compSet, ok = registry.CreateComponentSet(compType)
				if !ok {
					return nil, fmt.Errorf("ecs: component type %d not registered for group marshal", compType)
				}
			}
			if compSet == nil {
				continue
			}

			var csetData *SerializableSparseArrayData
			var err error
			if a != nil {
				csetData, err = compSet.marshalWithGroup(a, col)
			} else {
				cset, ok := compSet.(interface {
					Marshal() (*SerializableSparseArrayData, error)
				})
				if !ok {
					return nil, fmt.Errorf("ecs: component set type %d does not support serialization", compType)
				}
				csetData, err = cset.Marshal()
			}
			if err != nil {
				return nil, fmt.Errorf("ecs: marshal component %d: %w", compType, err)
			}
			data.ComponentSets = append(data.ComponentSets, SerializableComponentSetEntry{
				ComponentType: uint64(compType),
				Data:          *csetData,
			})
		}
	}

	// Serialize disposable types
	if len(w.disposableTypes) > 0 {
		data.DisposableTypes = make([]uint64, len(w.disposableTypes))
		for i, t := range w.disposableTypes {
			data.DisposableTypes[i] = uint64(t)
		}
	}

	// Serialize nomadic types
	if len(w.nomadicTypes) > 0 {
		data.NomadicTypes = make([]uint64, len(w.nomadicTypes))
		for i, t := range w.nomadicTypes {
			data.NomadicTypes[i] = uint64(t)
		}
	}

	return data, nil
}

// MarshalTo writes serializableWorld to a rockmem.Writer
func (w *serializableWorld) MarshalTo(writer rockmem.Writer) (uint, error) {
	data, err := w.Marshal()
	if err != nil {
		return 0, err
	}
	return data.WriteAsRoot(writer)
}

func (w *World) checkNoPendingForMarshal() error {
	if w.opLog != nil && w.opLog.HasPending() {
		return fmt.Errorf("ecs: cannot marshal world with pending component operations (call Update/Step first)")
	}
	w.destroyMu.Lock()
	hasDestroy := len(w.pendingDestroy) > 0
	w.destroyMu.Unlock()
	if hasDestroy {
		return fmt.Errorf("ecs: cannot marshal world with pending entity destructions (call Update/Step first)")
	}
	return nil
}

// MarshalTo writes world to a rockmem.Writer (convenience method)
// 只读联合导出组表列与 CSet 残段，不修改实时世界布局。
func (w *World) MarshalTo(writer rockmem.Writer) (uint, error) {
	data, err := w.Marshal()
	if err != nil {
		return 0, err
	}
	return data.WriteAsRoot(writer)
}

// Marshal serializes world to SerializableWorldData (convenience method)
// 只读联合导出组表列与 CSet 残段，不修改实时世界布局。
func (w *World) Marshal() (*SerializableWorldData, error) {
	if err := w.checkNoPendingForMarshal(); err != nil {
		return nil, err
	}
	return w.serializableWorld.marshalWithArchetypes(w.archetypes)
}

// Unmarshal deserializes SerializableWorldData to serializableWorld.
// R19：版本校验、重复组件池拒绝、组件逐个 error 传播（未注册类型报错而非静默跳过）、
// 任何失败原子返回（调用方负责在失败时不替换目标世界状态）
func (w *serializableWorld) Unmarshal(data *SerializableWorldData) error {
	if data.Version != SnapshotVersion {
		return fmt.Errorf("ecs: snapshot version mismatch: data=%d current=%d", data.Version, SnapshotVersion)
	}

	// 组件集先行恢复（最易失败），拒绝重复 ComponentType，成功后才提交到 world
	components := make(map[ComponentIntType]ComponentSet, len(data.ComponentSets))
	registry := GetComponentRegistry()
	for _, entry := range data.ComponentSets {
		compType := ComponentIntType(entry.ComponentType)
		if _, dup := components[compType]; dup {
			return fmt.Errorf("ecs: duplicate component set %d in snapshot", compType)
		}
		compSet, err := registry.UnmarshalComponentSet(compType, &entry.Data)
		if err != nil {
			return fmt.Errorf("ecs: unmarshal component %d: %w", compType, err)
		}
		components[compType] = compSet
	}

	entities := NewEntitySet()
	if err := entities.Unmarshal(&data.Entities); err != nil {
		return err
	}
	idGenerator := &EntityIDGenerator{}
	if err := idGenerator.Unmarshal(&data.IdGenerator); err != nil {
		return err
	}

	w.frame = data.Frame
	w.entities = entities
	w.idGenerator = idGenerator
	w.components = components

	// Deserialize disposable types
	if len(data.DisposableTypes) > 0 {
		w.disposableTypes = make([]ComponentIntType, len(data.DisposableTypes))
		for i, t := range data.DisposableTypes {
			w.disposableTypes[i] = ComponentIntType(t)
		}
	} else {
		w.disposableTypes = nil
	}

	// Deserialize nomadic types
	if len(data.NomadicTypes) > 0 {
		w.nomadicTypes = make([]ComponentIntType, len(data.NomadicTypes))
		for i, t := range data.NomadicTypes {
			w.nomadicTypes[i] = ComponentIntType(t)
		}
	} else {
		w.nomadicTypes = nil
	}
	return nil
}

// UnmarshalFrom reads serializableWorld from a rockmem.Reader
func (w *serializableWorld) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableWorldData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableWorldData root buffer")
	}
	data.ReadAsRoot(reader)
	return w.Unmarshal(data)
}

// Unmarshal deserializes SerializableWorldData to world (convenience method)。
// R19 运行时卫生：仅允许在初始化期（非 Running/Stop）的世界上原地恢复；
// 恢复整体替换组件表、清空旧组表行、丢弃遗留销毁队列与操作日志、递增组件版本使查询缓存失效。
// 失败原子返回，目标世界保持原状。
func (w *World) Unmarshal(data *SerializableWorldData) error {
	if w.status == WorldStatusRunning || w.status == WorldStatusStop {
		return fmt.Errorf("ecs: cannot unmarshal into world in status %v", w.status)
	}
	if err := w.serializableWorld.Unmarshal(data); err != nil {
		return err
	}
	if w.archetypes != nil {
		w.archetypes.clearRows()
	}
	// 丢弃快照前的运行时残留：遗留销毁/操作不得打进新世界
	w.destroyMu.Lock()
	w.pendingDestroy = nil
	w.destroyMu.Unlock()
	if w.opLog != nil {
		order := w.opLog.sysOrder
		w.opLog = NewOpLog(w, runtime.NumCPU())
		// 原位恢复保留已注册的系统。为它们重建空的专属队列，丢弃快照之前的命令。
		for _, id := range order {
			w.opLog.registerSystem(id)
		}
	}
	w.compVersion++ // 组件池整体替换：旧查询缓存全部失效
	// Set world pointer for all entities after deserialization
	w.entities.SetWorldForAll(w)
	return nil
}

// UnmarshalFrom reads world from a rockmem.Reader (convenience method)
func (w *World) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableWorldData{}
	if !data.ValidateAsRoot(reader) {
		return fmt.Errorf("ecs: invalid SerializableWorldData root buffer")
	}
	data.ReadAsRoot(reader)
	return w.Unmarshal(data)
}

// NewWorldFromData creates a new world from SerializableWorldData
// Note: This creates a minimal world without systems, optimizer, etc.
// The caller should configure these runtime components after deserialization.
func NewWorldFromData(data *SerializableWorldData, opts ...WorldOption) (*World, error) {
	c := &WorldConfig{}
	c.initDefault()
	for _, opt := range opts {
		opt(c)
	}

	w := &World{}
	w.status = WorldStatusInitializing
	w.config = c

	// Initialize runtime components (these are not serialized)
	flowExecOpt := WithFlowSyncMode()
	if c.ExecuteMode == ExecuteModeParallel {
		flowExecOpt = WithFlowASyncMode()
	}
	w.systems = newSystemFlow(w, flowExecOpt)
	w.opLog = NewOpLog(w, runtime.NumCPU())
	w.optimizer = newOptimizer(w)
	w.archetypes = NewArchetypeRegistry()

	// Unmarshal the serializable state
	if err := w.Unmarshal(data); err != nil {
		return nil, err
	}

	w.status = WorldStatusInitialized
	return w, nil
}
