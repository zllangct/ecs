package ecs

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// SnapshotVersion 快照格式版本（R19）：不兼容演进递增，恢复时校验拒绝
const SnapshotVersion uint64 = 1

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
// R19：校验快照与当前类型的布局一致性（EleSize 必须等于 sizeof(T)）与数据长度，
// 非法输入返回 error 且不修改目标（原子失败）
func (u *USet[T]) Unmarshal(data *SerializableUSetData) error {
	wantSize := uint64(TypeOf[T]().Size())
	if data.Len > 0 {
		if data.EleSize != wantSize {
			return fmt.Errorf("ecs: USet element size mismatch: snapshot=%d current=%d", data.EleSize, wantSize)
		}
		eleCount := int(data.Len)
		if eleCount < 0 || int64(eleCount) != data.Len {
			return fmt.Errorf("ecs: USet invalid length %d", data.Len)
		}
		byteSize := eleCount * int(data.EleSize)
		if byteSize < 0 || byteSize > len(data.Data) {
			return fmt.Errorf("ecs: USet data too short: need %d bytes, got %d", byteSize, len(data.Data))
		}

		u.data = u.a.alloc(eleCount, eleCount)
		dstPtr := unsafe.Pointer(&u.data[0])
		copy(unsafe.Slice((*byte)(dstPtr), byteSize), data.Data[:byteSize])
	} else {
		u.data = u.a.alloc(0, int(data.InitSize))
	}
	u.eleSize = data.EleSize
	u.len = data.Len
	u.initSize = data.InitSize
	return nil
}

// UnmarshalFrom reads USet[T] from a rockmem.Reader
func (u *USet[T]) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableUSetData{}
	data.ReadAsRoot(reader)
	return u.Unmarshal(data)
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
func (s *SparseArray[K, V]) Unmarshal(data *SerializableSparseArrayData) error {
	// Unmarshal the embedded USet
	if err := s.USet.Unmarshal(&data.USetData); err != nil {
		return err
	}

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
	return nil
}

// UnmarshalFrom reads SparseArray[K, V] from a rockmem.Reader
func (s *SparseArray[K, V]) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableSparseArrayData{}
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
	data := c.SparseArray.Marshal()

	var probe TP
	if _, ok := any(probe).(rockmemMessage); !ok || c.Len() == 0 {
		return data, nil
	}

	n := c.Len()
	offsets := make([]int64, n+1)
	var blob []byte
	w := rockmem.NewWriter()
	for i := 0; i < n; i++ {
		w.Reset()
		m := any(TP(&c.data[i])).(rockmemMessage)
		if _, err := m.WriteAsRoot(w); err != nil {
			return nil, fmt.Errorf("ecs: marshal component %d: %w", i, err)
		}
		offsets[i] = int64(len(blob))
		blob = append(blob, w.Bytes()...)
	}
	offsets[n] = int64(len(blob))

	hdr := make([]byte, (n+1)*8)
	for i, o := range offsets {
		binary.LittleEndian.PutUint64(hdr[i*8:], uint64(o))
	}
	data.USetData.Data = append(hdr, blob...)
	data.USetData.EleSize = 0 // IDL 路径标记
	return data, nil
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
// R19：校验偏移表单调且在数据界内，损坏快照返回 error（原子失败）
func (c *CSet[T, TP]) unmarshalIDL(data *SerializableSparseArrayData) error {
	raw := data.USetData.Data
	n := int(data.USetData.Len)
	hdrLen := int64((n + 1) * 8)
	if int64(len(raw)) < hdrLen {
		return fmt.Errorf("ecs: IDL snapshot header truncated: need %d, got %d", hdrLen, len(raw))
	}
	offsets := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		offsets[i] = int64(binary.LittleEndian.Uint64(raw[i*8:]))
		if i > 0 && offsets[i] < offsets[i-1] {
			return fmt.Errorf("ecs: IDL offset table not monotonic at %d", i)
		}
	}
	blobLen := int64(len(raw)) - hdrLen
	if offsets[n] > blobLen {
		return fmt.Errorf("ecs: IDL blob truncated: need %d, got %d", offsets[n], blobLen)
	}

	c.USet.eleSize = uint64(TypeOf[T]().Size())
	c.USet.len = int64(n)
	c.USet.initSize = data.USetData.InitSize
	c.USet.data = c.USet.a.alloc(n, n)
	for i := 0; i < n; i++ {
		blob := raw[hdrLen+offsets[i] : hdrLen+offsets[i+1]]
		m := any(TP(&c.USet.data[i])).(rockmemMessage)
		m.ReadAsRoot(rockmem.NewReader(blob))
	}

	// SparseArray 元数据（与 SparseArray.Unmarshal 相同部分）
	c.maxKey = EntityIndex(data.MaxKey)
	c.shrinkThreshold = data.ShrinkThreshold
	c.initSize = int(data.InitSize)
	c.isKOrder = data.IsKOrder
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
	count := len(data.EntityIds)
	if count > 0 {
		if len(data.CompoundOffsets) != count || len(data.CompoundLengths) != count {
			return fmt.Errorf("ecs: entity parallel arrays length mismatch: ids=%d offsets=%d lengths=%d",
				count, len(data.CompoundOffsets), len(data.CompoundLengths))
		}
		prevEnd := int32(0)
		for i := 0; i < count; i++ {
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

	// Deserialize all EntityInfos from parallel arrays
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
	return nil
}

// UnmarshalFrom reads EntitySet from a rockmem.Reader
func (es *EntitySet) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableEntitySetData{}
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
// R19：校验 freelist 不变量（pending 不越界），非法返回 error
func (e *EntityIDGenerator) Unmarshal(data *SerializableEntityIDGeneratorData) error {
	if data.Pending < 0 || int64(len(data.Ids)) < int64(data.Pending) {
		return fmt.Errorf("ecs: id generator pending %d exceeds ids %d", data.Pending, len(data.Ids))
	}
	e.free = EntityIndex(data.Free)
	e.pending = EntityIndex(data.Pending)
	e.len = data.Len
	e.delayFree = data.DelayFree
	e.delayCap = data.DelayCap

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

	if e.delayCap <= 0 {
		e.delayCap = 10 // 与 NewEntityIDGenerator 默认值对齐
	}

	// Convert int64 slice to removeDelay slice
	if len(data.RemoveDelay) > 0 {
		e.removeDelay = make([]ReuseID, len(data.RemoveDelay))
		for i, id := range data.RemoveDelay {
			e.removeDelay[i] = Entity(id).toReuseID()
		}
		if int32(len(e.removeDelay)) < e.delayCap {
			e.removeDelay = append(e.removeDelay, make([]ReuseID, int(e.delayCap)-len(e.removeDelay))...)
		}
	} else {
		e.removeDelay = make([]ReuseID, e.delayCap)
	}
	return nil
}

// UnmarshalFrom reads EntityIDGenerator from a rockmem.Reader
func (e *EntityIDGenerator) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableEntityIDGeneratorData{}
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

	// Serialize component sets
	if len(w.components) > 0 {
		data.ComponentSets = make([]SerializableComponentSetEntry, 0, len(w.components))
		for compType, compSet := range w.components {
			// Cast to CSet to access Marshal method
			cset, ok := compSet.(interface {
				Marshal() (*SerializableSparseArrayData, error)
			})
			if !ok {
				return nil, fmt.Errorf("ecs: component set type %d does not support serialization", compType)
			}
			csetData, err := cset.Marshal()
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

// MarshalTo writes world to a rockmem.Writer (convenience method)
// 组表数据先散回 CSet 导出（格式兼容），完成后重新收拢恢复布局。
func (w *World) MarshalTo(writer rockmem.Writer) (uint, error) {
	w.scatterAll()
	defer w.absorbAll()
	return w.serializableWorld.MarshalTo(writer)
}

// Marshal serializes world to SerializableWorldData (convenience method)
// 组表数据先散回 CSet 导出（格式兼容），完成后重新收拢恢复布局。
func (w *World) Marshal() (*SerializableWorldData, error) {
	w.scatterAll()
	defer w.absorbAll()
	return w.serializableWorld.Marshal()
}

// Unmarshal deserializes SerializableWorldData to serializableWorld.
// R19：版本校验、组件逐个 error 传播（未注册类型报错而非静默跳过）、
// 任何失败原子返回（调用方负责在失败时不替换目标世界状态）
func (w *serializableWorld) Unmarshal(data *SerializableWorldData) error {
	if data.Version != SnapshotVersion {
		return fmt.Errorf("ecs: snapshot version mismatch: data=%d current=%d", data.Version, SnapshotVersion)
	}

	// 组件集先行恢复（最易失败），成功后才提交到 world
	components := make(map[ComponentIntType]ComponentSet, len(data.ComponentSets))
	registry := GetComponentRegistry()
	for _, entry := range data.ComponentSets {
		compType := ComponentIntType(entry.ComponentType)
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
	data.ReadAsRoot(reader)
	return w.Unmarshal(data)
}

// Unmarshal deserializes SerializableWorldData to world (convenience method)。
// R19 运行时卫生：仅允许在停止/初始化完成（非 Running）的世界上原地恢复；
// 恢复整体替换组件表、丢弃遗留销毁队列与操作日志、递增组件版本使查询缓存失效。
// 失败原子返回，目标世界保持原状。
func (w *World) Unmarshal(data *SerializableWorldData) error {
	if w.status == WorldStatusRunning {
		return fmt.Errorf("ecs: cannot unmarshal into a running world (quiesce first)")
	}
	if err := w.serializableWorld.Unmarshal(data); err != nil {
		return err
	}
	// 丢弃快照前的运行时残留：遗留销毁/操作不得打进新世界
	w.pendingDestroy = nil
	if w.opLog != nil {
		w.opLog = NewOpLog(w, runtime.NumCPU())
	}
	w.compVersion++ // 组件池整体替换：旧查询缓存全部失效
	// Set world pointer for all entities after deserialization
	w.entities.SetWorldForAll(w)
	return nil
}

// UnmarshalFrom reads world from a rockmem.Reader (convenience method)
func (w *World) UnmarshalFrom(reader *rockmem.Reader) error {
	data := &SerializableWorldData{}
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
