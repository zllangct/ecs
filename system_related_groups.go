package ecs

import (
	"iter"
	"sort"
)

// systemTreeNode system tree systemTreeNode
type systemTreeNode struct {
	parent   *systemTreeNode
	children []*systemTreeNode
	val      SystemInfo
}

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

func (p *systemTreeNode) attach(node *systemTreeNode) {
	isAttached := false
	for i := 0; i < len(p.children); i++ {
		if p.children[i].isFriend(node) {
			p.children[i].attach(node)
			isAttached = true
			break
		}
	}
	if !isAttached {
		if p.val == node.val {
			return
		}
		p.children = append(p.children, node)
	}
}

// SystemRelatedGroups system group ordered by interrelation
type SystemRelatedGroups struct {
	systems      []*systemTreeNode
	ref          map[ComponentIntType]int
	root         *systemTreeNode
	order        Order
	batchTotal   int
	maxPeerBatch int
	ordered      bool
}

func NewSystemRelatedGroups() SystemTraverser {
	return &SystemRelatedGroups{
		systems: make([]*systemTreeNode, 0),
		ref:     map[ComponentIntType]int{},
		ordered: true,
		root: &systemTreeNode{
			parent:   nil,
			children: []*systemTreeNode{},
			val:      nil,
		},
	}
}

func (s *SystemRelatedGroups) refCount(rqs []Dependency) int {
	ref := 0
	for _, com := range rqs {
		ref += s.ref[com.intType()] - 1
	}
	return ref
}

func (s *SystemRelatedGroups) resort() {
	if s.ordered {
		return
	}
	sort.Slice(s.systems, func(i, j int) bool {
		return s.refCount(s.systems[i].val.getDeps()) >
			s.refCount(s.systems[j].val.getDeps())
	})

	s.root.children = []*systemTreeNode{}
	for _, node := range s.systems {
		node.children = []*systemTreeNode{}
		s.root.attach(node)
	}
	s.ordered = true

	s.batchTotal = 0
	s.maxPeerBatch = 0

	var top []*systemTreeNode = s.root.children
	for len(top) > 0 {
		count := 0
		temp := top
		top = make([]*systemTreeNode, 0)
		for _, node := range temp {
			count++
			top = append(top, node.children...)
		}
		if count > s.maxPeerBatch {
			s.maxPeerBatch = count
		}
		s.batchTotal++
	}
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

// get all systems
func (s *SystemRelatedGroups) all() []SystemInfo {
	systems := make([]SystemInfo, len(s.systems))
	for i, n := range s.systems {
		systems[i] = n.val
	}
	return systems
}

// insert system
func (s *SystemRelatedGroups) add(sys SystemInfo) error {
	//get system's required components
	rqs := sys.getDeps()
	if len(rqs) == 0 {
		return nil
	}
	//reference count
	for _, com := range rqs {
		if _, ok := s.ref[com.intType()]; ok {
			s.ref[com.intType()] += 1
		} else {
			s.ref[com.intType()] = 1
		}
	}
	//add system
	node := &systemTreeNode{
		children: make([]*systemTreeNode, 0),
		val:      sys,
	}
	s.systems = append(s.systems, node)
	//set unordered
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
	//get system's required components
	rqs := sys.getDeps()
	has := false
	for i, system := range s.systems {
		if system.val.id() == sys.id() {
			s.systems = append(s.systems[:i], s.systems[i+1:]...)
			has = true
			break
		}
	}
	if !has {
		return nil
	}
	//reference count
	for _, com := range rqs {
		if _, ok := s.ref[com.intType()]; ok {
			s.ref[com.intType()] -= 1
		} else {
			panic("component ref wrong")
		}
	}
	//set unordered
	s.ordered = false
	return nil
}

func (s *SystemRelatedGroups) independentGroups() iter.Seq[[]SystemInfo] {
	if !s.ordered {
		s.resort()
	}

	if s.maxPeerBatch == 0 {
		return func(yield func([]SystemInfo) bool) {}
	}
	top := make([]*systemTreeNode, s.maxPeerBatch)
	topTemp := make([]*systemTreeNode, s.maxPeerBatch)
	buffer := make([]SystemInfo, s.maxPeerBatch)
	size := 0
	copy(top, s.root.children)
	topSize := len(s.root.children)

	return func(yield func([]SystemInfo) bool) {
		for {
			if topSize == 0 {
				size = 0
				return
			}
			topTemp, top = top, topTemp
			tempSize := topSize
			topSize = 0
			size = 0
			for i := 0; i < tempSize; i++ {
				buffer[size] = topTemp[i].val
				size++
				for j := 0; j < len(topTemp[i].children); j++ {
					top[topSize+j] = topTemp[i].children[j]
				}
				topSize += len(topTemp[i].children)
			}

			if size == 0 {
				return
			}
			if !yield(buffer[:size]) {
				return
			}
		}
	}
}
