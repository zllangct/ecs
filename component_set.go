package ecs

import (
	"iter"
	"slices"
	"unsafe"
)

type ComponentSet interface {
	Len() int
	Reset()
	Free()
	Add(entity Entity, comp Component) Component
	// AddRaw 从裸内存拷贝一个元素入集（组表拆行散回用，src 指向组件定长内存）
	AddRaw(entity Entity, src unsafe.Pointer)
	// Exist 实体是否有组件（稀疏表一次索引，Query 过滤用）
	Exist(index EntityIndex) bool
	// IsKeyOrdered 密集数组是否按 key 升序（merge-join 查询路径前提）
	IsKeyOrdered() bool
	Get(entity Entity) Component
	Remove(entity Entity)
	RemoveAndReturn(entity Entity) Component
	EntityIndexes() []EntityIndex
	get(index EntityIndex) unsafe.Pointer
	// iterPtr 以 unsafe.Pointer 遍历密集数组（零拷贝），供只读视图等
	// 不持有具体类型参数的访问路径使用
	iterPtr() iter.Seq2[EntityIndex, unsafe.Pointer]
	Sort()
	// marshalWithGroup 只读联合导出组表列与 CSet 残段，不修改实时世界布局
	marshalWithGroup(a *Archetype, col int) (*SerializableSparseArrayData, error)
}

type CSet[T any, TP ComponentPointer[T]] struct {
	SparseArray[EntityIndex, T]
}

// 用于识别泛型参数未知的原生池。嵌入 CSet 的自定义池返回其内层池。
func (c *CSet[T, TP]) nativeCSet() ComponentSet { return c }

func NewCSet[T any, TP ComponentPointer[T]](initSize ...int) *CSet[T, TP] {
	c := &CSet[T, TP]{
		SparseArray: *NewSparseArray[EntityIndex, T](initSize...),
	}
	return c
}

func (c *CSet[T, TP]) EntityIndexes() []EntityIndex {
	return *(*[]EntityIndex)(unsafe.Pointer(&c.idx2Key))
}

func (c *CSet[T, TP]) Add(entity Entity, comp Component) Component {
	if comp == nil {
		return nil
	}
	ptr := (*iface)(unsafe.Pointer(&comp)).data
	if ptr == nil {
		return nil
	}
	index := entity.Index()
	// 注意：此处复刻 Go runtime 的 iface 布局直接取 data 指针（见 internal_type_mock.go），
	// 零断言开销，但与 Go 版本的接口内存布局绑定，升级 toolchain 需回归验证。
	data := c.SparseArray.Add(index, (*T)(ptr))
	if data == nil {
		return nil
	}
	return any(data).(Component)
}

func (c *CSet[T, TP]) Free() {
	c.SparseArray.Free()
}

func (c *CSet[T, TP]) marshalWithGroup(a *Archetype, col int) (*SerializableSparseArrayData, error) {
	source := c
	if a.Len() > 0 {
		// 行序也是存档数据。恢复时移出共同尾部的组行，保留前缀中各自的残段顺序。
		source = NewCSet[T, TP](c.Len() + a.Len())
		for idx, val := range c.Iter() {
			source.SparseArray.Add(idx, val)
		}
		for row, idx := range a.entities {
			source.SparseArray.Add(idx, (*T)(a.cellPtr(col, row)))
		}
	}
	data, err := source.Marshal()
	if err != nil {
		return nil, err
	}
	// 还原真实组件池的配置。稀疏索引的预留尾部及失序标志的历史值均不影响
	// 元素顺序，将它们规范化，确保组行全部拆出后的下一份快照仍可逐字节复现。
	data.InitSize = int32(c.SparseArray.initSize)
	data.USetData.InitSize = c.USet.initSize
	data.ShrinkThreshold = c.shrinkThreshold
	end := len(data.Indices)
	for end > 0 && data.Indices[end-1] == 0 {
		end--
	}
	data.Indices = data.Indices[:end]
	data.IsKOrder = data.IsKOrder || slices.IsSorted(data.Idx2Key)
	return data, nil
}

// AddRaw 从裸内存拷贝一个元素入集（USet.Add 为拷贝语义，不与 src 共享内存）
func (c *CSet[T, TP]) AddRaw(entity Entity, src unsafe.Pointer) {
	c.SparseArray.Add(entity.Index(), (*T)(src))
}

// Exist 实体是否有组件（稀疏表一次索引）
func (c *CSet[T, TP]) Exist(index EntityIndex) bool {
	return c.SparseArray.Exist(index)
}

// IsKeyOrdered 密集数组是否按 key 升序
func (c *CSet[T, TP]) IsKeyOrdered() bool {
	return c.SparseArray.isKOrder
}

func (c *CSet[T, TP]) Remove(entity Entity) {
	key := entity.Index()
	if idx := c.denseIndex(key); idx >= 0 {
		c.removeAt(key, idx)
	}
}

func (c *CSet[T, TP]) RemoveAndReturn(entity Entity) Component {
	r := c.SparseArray.Remove(entity.Index())
	if r == nil {
		return nil
	}
	return TP(r)
}

func (c *CSet[T, TP]) Get(entity Entity) Component {
	data := c.SparseArray.Get(entity.Index())
	if data == nil {
		return nil
	}
	return any(data).(Component)
}

func (c *CSet[T, TP]) iterPtr() iter.Seq2[EntityIndex, unsafe.Pointer] {
	return func(yield func(EntityIndex, unsafe.Pointer) bool) {
		for index, p := range c.Iter() {
			if !yield(index, unsafe.Pointer(p)) {
				return
			}
		}
	}
}

func (c *CSet[T, TP]) get(entityIndex EntityIndex) unsafe.Pointer {
	data := c.SparseArray.Get(entityIndex)
	if data == nil {
		return nil
	}
	return unsafe.Pointer(data)
}

func (c *CSet[T, TP]) getByIndex(index EntityIndex) *T {
	return c.SparseArray.Get(index)
}
