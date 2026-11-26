package ecs

import (
	"iter"
)

// systemTreeNode system tree systemTreeNode
type systemTreeNode struct {
	parent   *systemTreeNode
	children []*systemTreeNode
	val      SystemInfo
}

// isFriend 判断两个 system 是否因共享组件依赖而"相关"（任一方可写即冲突）。
// 命名取"好友"形象义：好友必须排进不同 batch 串行，不能并行。
func (p *systemTreeNode) isFriend(node *systemTreeNode) bool {
	for _, r := range p.val.getDeps() {
		for _, rTarget := range node.val.getDeps() {
			if r.intType() == rTarget.intType() {
				if r.readonly() && rTarget.readonly() {
					continue
				}
				return true
			}
		}
	}
	return false
}

// SystemRelatedGroups 按冲突关系与注册顺序分层调度系统。
// 同 Order 档位内：互不冲突（含无依赖/全只读共享）的系统进入同一 batch 并行；
// 存在读写冲突的系统严格按注册顺序分入后续 batch 串行。
type SystemRelatedGroups struct {
	systems      []*systemTreeNode
	batches      [][]SystemInfo
	order        Order
	batchTotal   int
	maxPeerBatch int
	ordered      bool
}

func NewSystemRelatedGroups() SystemTraverser {
	return &SystemRelatedGroups{
		systems: make([]*systemTreeNode, 0),
		ordered: true,
	}
}

func (s *SystemRelatedGroups) resort() {
	if s.ordered {
		return
	}
	n := len(s.systems)
	if n == 0 {
		s.batches = nil
		s.batchTotal = 0
		s.maxPeerBatch = 0
		s.ordered = true
		return
	}

	// 按注册顺序逐个分配最早无冲突批次：
	// 若 i < j 且 systems[i] 与 systems[j] 冲突，则 layer[j] >= layer[i] + 1，
	// 既保证任意同批内全对全无读写冲突，又保证同档位冲突系统严格按注册顺序执行。
	layers := make([]int, n)
	maxLayer := 0
	for j := 0; j < n; j++ {
		layer := 0
		for i := 0; i < j; i++ {
			if s.systems[i].isFriend(s.systems[j]) && layers[i]+1 > layer {
				layer = layers[i] + 1
			}
		}
		layers[j] = layer
		if layer > maxLayer {
			maxLayer = layer
		}
	}

	s.batches = make([][]SystemInfo, maxLayer+1)
	for i, node := range s.systems {
		l := layers[i]
		s.batches[l] = append(s.batches[l], node.val)
	}
	s.batchTotal = len(s.batches)
	s.maxPeerBatch = 0
	for _, b := range s.batches {
		if len(b) > s.maxPeerBatch {
			s.maxPeerBatch = len(b)
		}
	}
	s.ordered = true
}

func (s *SystemRelatedGroups) systemCount() int {
	return len(s.systems)
}

func (s *SystemRelatedGroups) getBatchCount() int {
	return s.batchTotal
}

func (s *SystemRelatedGroups) getMaxCountPeerBatch() int {
	return s.maxPeerBatch
}

func (s *SystemRelatedGroups) count() int {
	return len(s.systems)
}

func (s *SystemRelatedGroups) getOrder() Order {
	return s.order
}

func (s *SystemRelatedGroups) setOrder(order Order) {
	s.order = order
}

// get all systems (preserving registration order)
func (s *SystemRelatedGroups) all() []SystemInfo {
	systems := make([]SystemInfo, len(s.systems))
	for i, n := range s.systems {
		systems[i] = n.val
	}
	return systems
}

// insert system (including zero-dependency systems)
func (s *SystemRelatedGroups) add(sys SystemInfo) error {
	node := &systemTreeNode{
		val: sys,
	}
	s.systems = append(s.systems, node)
	s.ordered = false
	return nil
}

// has system
func (s *SystemRelatedGroups) has(sys SystemInfo) bool {
	for _, system := range s.systems {
		if system.val.id() == sys.id() {
			return true
		}
	}
	return false
}

// remove system
func (s *SystemRelatedGroups) remove(sys SystemInfo) error {
	for i, system := range s.systems {
		if system.val.id() == sys.id() {
			s.systems = append(s.systems[:i], s.systems[i+1:]...)
			s.ordered = false
			return nil
		}
	}
	return nil
}

func (s *SystemRelatedGroups) independentGroups() iter.Seq[[]SystemInfo] {
	if !s.ordered {
		s.resort()
	}
	if len(s.batches) == 0 {
		return func(yield func([]SystemInfo) bool) {}
	}
	batches := s.batches
	return func(yield func([]SystemInfo) bool) {
		for _, batch := range batches {
			if len(batch) == 0 {
				continue
			}
			if !yield(batch) {
				return
			}
		}
	}
}
