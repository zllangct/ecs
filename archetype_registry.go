package ecs

import (
	"fmt"
	"sort"
	"unsafe"
)

// ArchetypeRegistry 注册期收集 System 声明的固定访问组 → Finalize 重叠合并 → 建表。
// Finalize 后只读（ownerOf O(1) 分派，执行期无锁访问）。
type ArchetypeRegistry struct {
	declared [][]ComponentIntType
	sizes    map[ComponentIntType]int
	groups   map[FixedCompound]*Archetype
	ownerOf  map[ComponentIntType]*Archetype
	sealed   bool
}

func NewArchetypeRegistry() *ArchetypeRegistry {
	return &ArchetypeRegistry{
		sizes:   map[ComponentIntType]int{},
		groups:  map[FixedCompound]*Archetype{},
		ownerOf: map[ComponentIntType]*Archetype{},
	}
}

// Declare 声明一个固定访问组（types 不要求有序，至少 2 个类型）。
// 仅 Finalize 前可调用（注册期）。
func (r *ArchetypeRegistry) Declare(types []ComponentIntType, sizes map[ComponentIntType]int) error {
	if r.sealed {
		return fmt.Errorf("ecs: archetype registry sealed")
	}
	if len(types) < 2 {
		return fmt.Errorf("ecs: group needs >= 2 component types, got %v", types)
	}
	cp := make([]ComponentIntType, len(types))
	copy(cp, types)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	r.declared = append(r.declared, cp)
	for _, it := range cp {
		s, ok := sizes[it]
		if !ok || s <= 0 {
			return fmt.Errorf("ecs: component type %d size unknown (not registered?)", it)
		}
		r.sizes[it] = s
	}
	return nil
}

// Finalize 将重叠声明按传递闭包合并为并集组，建表并生成 ownerOf。
func (r *ArchetypeRegistry) Finalize() error {
	if r.sealed {
		return nil
	}
	// 并查集：含共同组件类型的声明合并
	parent := make([]int, len(r.declared))
	for i := range parent {
		parent[i] = i
	}
	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	seen := map[ComponentIntType]int{} // type → 首次出现的声明下标
	for i, decl := range r.declared {
		for _, it := range decl {
			if j, ok := seen[it]; ok {
				pi, pj := find(i), find(j)
				if pi != pj {
					parent[pi] = pj
				}
			} else {
				seen[it] = i
			}
		}
	}
	merged := map[int]map[ComponentIntType]struct{}{}
	for i, decl := range r.declared {
		root := find(i)
		if merged[root] == nil {
			merged[root] = map[ComponentIntType]struct{}{}
		}
		for _, it := range decl {
			merged[root][it] = struct{}{}
		}
	}
	for _, set := range merged {
		types := make([]ComponentIntType, 0, len(set))
		for it := range set {
			types = append(types, it)
		}
		sort.Slice(types, func(i, j int) bool { return types[i] < types[j] })
		sizes := make([]int, len(types))
		for i, it := range types {
			sizes[i] = r.sizes[it]
		}
		a := NewArchetype(types, sizes)
		key := NewFixedCompound(a.compound())
		r.groups[key] = a
		for _, it := range types {
			r.ownerOf[it] = a
		}
	}
	r.sealed = true
	return nil
}

// owner 组件类型所属组表
func (r *ArchetypeRegistry) owner(it ComponentIntType) (*Archetype, bool) {
	a, ok := r.ownerOf[it]
	return a, ok
}

// scatterAll 把所有组表行散回 CSet（序列化导出前调用，世界须处于静止点）。
// 散回后世界语义不变（组表是纯布局层）。
func (w *World) scatterAll() {
	for _, a := range w.archetypes.groups {
		for a.Len() > 0 {
			index := a.entities[0]
			info := w.entities.getByIndex(index)
			if info == nil {
				a.removeRow(index)
				continue
			}
			for col, it := range a.types {
				set, ok := w.getComponentSet(it)
				if !ok {
					continue
				}
				set.AddRaw(info.entity, a.cellPtr(col, 0))
			}
			a.removeRow(index) // swap-remove 后 0 号行是新行，原地继续
		}
	}
}

// absorbAll 按 ownerOf 把 CSet 中集齐组组件的实体收拢入行
// （Unmarshal 恢复后、Marshal 导出后调用）。
func (w *World) absorbAll() {
	for _, a := range w.archetypes.groups {
		// 候选 = 组内最小 CSet 的实体集合
		var minSet ComponentSet
		for _, it := range a.types {
			if s, ok := w.getComponentSet(it); ok && (minSet == nil || s.Len() < minSet.Len()) {
				minSet = s
			}
		}
		if minSet == nil {
			continue
		}
		touched := make(map[EntityIndex]struct{}, minSet.Len())
		for _, index := range minSet.EntityIndexes() {
			touched[index] = struct{}{}
		}
		w.syncGroup(a, touched)
	}
}

// syncGroup 帧同步点组成员资格重评估（两阶段 flush 的第二相，单线程/组间并发安全——
// 各组类型集合互不重叠，实体行只属一张表）。
// 对触及实体：在表中但已失去组内组件 → 拆行散回 CSet；不在表中但已集齐 → 收拢入行。
func (w *World) syncGroup(a *Archetype, entities map[EntityIndex]struct{}) {
	for index := range entities {
		info := w.entities.getByIndex(index)
		if info == nil {
			continue // 已销毁
		}
		_, inTable := a.rowOf(index)
		full := true
		for _, it := range a.types {
			if !info.compound.Exist(it) {
				full = false
				break
			}
		}
		switch {
		case inTable && !full:
			// 拆行：仍持有的组内组件散回各自 CSet，再删行
			row, _ := a.rowOf(index)
			for col, it := range a.types {
				if !info.compound.Exist(it) {
					continue // 已删类型无数据需散
				}
				set, ok := w.getComponentSet(it)
				if !ok {
					continue
				}
				set.AddRaw(info.entity, a.cellPtr(col, row))
			}
			a.removeRow(index)
		case !inTable && full:
			// 收拢：从各 CSet 拷贝入行并从 CSet 删除
			row := a.addRow(index)
			for col, it := range a.types {
				set, ok := w.getComponentSet(it)
				if !ok {
					continue
				}
				p := set.get(index)
				if p == nil {
					continue
				}
				size := a.sizes[col]
				src := unsafe.Slice((*byte)(p), size)
				dst := unsafe.Slice((*byte)(a.cellPtr(col, row)), size)
				copy(dst, src)
				set.Remove(info.entity)
			}
		}
	}
}
