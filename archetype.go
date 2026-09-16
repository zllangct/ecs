package ecs

import (
	"sort"
	"unsafe"
)

// Archetype 组表：一个固定访问组（组件类型集合 G）对应一张表。
// 不变量：实体拥有 G 全集 ⟺ 在表中占一行；组件实例只存一处（表列或 CSet 残段）。
// 仅在帧同步点被单线程修改；system 执行期只读。
type Archetype struct {
	types    []ComponentIntType // 升序
	sizes    []int              // 与 types 平行，列 stride
	columns  [][]byte           // 与 types 平行，len = rows*size
	entities []EntityIndex      // 行 → 实体
	rowIdx   []int32            // 稀疏: EntityIndex → 行号+1，0=不在表（仿 SparseArray.indices）
}

// NewArchetype types/sizes 为平行数组；内部按 types 升序重排（同步重排 sizes）。
func NewArchetype(types []ComponentIntType, sizes []int) *Archetype {
	idx := make([]int, len(types))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return types[idx[a]] < types[idx[b]] })
	st := make([]ComponentIntType, len(types))
	ss := make([]int, len(sizes))
	for i, j := range idx {
		st[i] = types[j]
		ss[i] = sizes[j]
	}
	return &Archetype{
		types:   st,
		sizes:   ss,
		columns: make([][]byte, len(st)),
	}
}

func (a *Archetype) Len() int { return len(a.entities) }

// has 组内是否包含组件类型
func (a *Archetype) has(it ComponentIntType) bool {
	i := sort.Search(len(a.types), func(i int) bool { return a.types[i] >= it })
	return i < len(a.types) && a.types[i] == it
}

// colOf 组件类型的列号
func (a *Archetype) colOf(it ComponentIntType) (int, bool) {
	i := sort.Search(len(a.types), func(i int) bool { return a.types[i] >= it })
	if i < len(a.types) && a.types[i] == it {
		return i, true
	}
	return -1, false
}

// rowOf 实体所在行
func (a *Archetype) rowOf(index EntityIndex) (int, bool) {
	if index < 0 || int(index) >= len(a.rowIdx) {
		return -1, false
	}
	r := a.rowIdx[index] - 1
	if r < 0 {
		return -1, false
	}
	return int(r), true
}

// addRow 追加一行（各列按 stride 扩零值），返回行号。调用方保证实体不在表中。
func (a *Archetype) addRow(index EntityIndex) int {
	row := len(a.entities)
	a.entities = append(a.entities, index)
	for col, size := range a.sizes {
		a.columns[col] = append(a.columns[col], make([]byte, size)...)
	}
	if int(index) >= len(a.rowIdx) {
		n := make([]int32, int(index)*2+1)
		copy(n, a.rowIdx)
		a.rowIdx = n
	}
	a.rowIdx[index] = int32(row + 1)
	return row
}

// removeRow swap-remove：末行顶替被删行。返回被顶落实体（无则 -1）。
func (a *Archetype) removeRow(index EntityIndex) EntityIndex {
	row, ok := a.rowOf(index)
	if !ok {
		return -1
	}
	last := len(a.entities) - 1
	lastEntity := a.entities[last]
	if row != last {
		for col, size := range a.sizes {
			src := a.columns[col][last*size : (last+1)*size]
			dst := a.columns[col][row*size : (row+1)*size]
			copy(dst, src)
		}
		a.entities[row] = lastEntity
		a.rowIdx[lastEntity] = int32(row + 1)
	}
	a.entities = a.entities[:last]
	for col, size := range a.sizes {
		a.columns[col] = a.columns[col][:last*size]
	}
	a.rowIdx[index] = 0
	return lastEntity
}

// cellPtr 列 col 行 row 的元素指针
func (a *Archetype) cellPtr(col, row int) unsafe.Pointer {
	size := a.sizes[col]
	return unsafe.Pointer(&a.columns[col][row*size])
}

// entityIndexes 行序实体索引（Query 主迭代用）
func (a *Archetype) entityIndexes() []EntityIndex { return a.entities }

// compound 组身份
func (a *Archetype) compound() Compound {
	c := NewCompound(len(a.types))
	for _, it := range a.types {
		c.Add(it)
	}
	return c
}
