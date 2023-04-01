package ecs

import (
	"reflect"
	"sort"
)

var emptySystemGroupIterator = &SystemGroupIterator{}

// Node system tree node
type Node struct {
	parent   *Node
	children []*Node
	val      ISystem
}

func (p *Node) isFriend(node *Node) bool {
	for com, r := range p.val.GetRequirements() {
		for comTarget, rTarget := range node.val.GetRequirements() {
			if comTarget == com {
				if r.getPermission() == ComponentReadOnly && rTarget.getPermission() == ComponentReadOnly {
					continue
				}
				return true
			}
		}
	}
	return false
}

func (p *Node) attach(node *Node) {
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
			Log.Error("repeated system")
			return
		}
		p.children = append(p.children, node)
	}
}

// SystemGroup system group ordered by interrelation
type SystemGroup struct {
	SystemGroupIterator
	systems      []*Node
	ref          map[reflect.Type]int
	root         *Node
	order        Order
	batchTotal   int
	maxPeerBatch int
	ordered      bool
}

func NewSystemGroup() *SystemGroup {
	return &SystemGroup{
		systems: make([]*Node, 0),
		ref:     map[reflect.Type]int{},
		ordered: true,
		root: &Node{
			parent:   nil,
			children: []*Node{},
			val:      nil,
		},
	}
}

func (p *SystemGroup) refCount(rqs map[reflect.Type]IRequirement) int {
	ref := 0
	for com, _ := range rqs {
		ref += p.ref[com] - 1
	}
	return ref
}

func (p *SystemGroup) resort() {
	if p.ordered {
		return
	}
	sort.Slice(p.systems, func(i, j int) bool {
		return p.refCount(p.systems[i].val.GetRequirements()) >
			p.refCount(p.systems[j].val.GetRequirements())
	})

	p.root.children = []*Node{}
	for _, node := range p.systems {
		node.children = []*Node{}
		p.root.attach(node)
	}
	p.ordered = true

	p.batchTotal = 0
	p.maxPeerBatch = 0

	var top []*Node = p.root.children
	for len(top) > 0 {
		count := 0
		temp := top
		top = make([]*Node, 0)
		for _, node := range temp {
			count++
			top = append(top, node.children...)
		}
		if count > p.maxPeerBatch {
			p.maxPeerBatch = count
		}
		p.batchTotal++
	}

	p.resetIter()
}

func (p *SystemGroup) resetIter() {
	if p.group == nil {
		p.group = p
	}
	curLen := len(p.SystemGroupIterator.top)
	if curLen-p.maxPeerBatch == 1 {
		p.SystemGroupIterator.top = p.SystemGroupIterator.top[:p.maxPeerBatch]
		p.SystemGroupIterator.topTemp = p.SystemGroupIterator.topTemp[:p.maxPeerBatch]
		p.SystemGroupIterator.buffer = p.SystemGroupIterator.buffer[:p.maxPeerBatch]
	} else if curLen-p.maxPeerBatch == -1 {
		p.SystemGroupIterator.top = append(p.SystemGroupIterator.top, (*Node)(nil))
		p.SystemGroupIterator.topTemp = append(p.SystemGroupIterator.topTemp, (*Node)(nil))
		p.SystemGroupIterator.buffer = append(p.SystemGroupIterator.buffer, ISystem(nil))
	} else if curLen-p.maxPeerBatch == 0 {
		// do nothing
	} else {
		p.SystemGroupIterator.top = make([]*Node, p.maxPeerBatch)
		p.SystemGroupIterator.topTemp = make([]*Node, p.maxPeerBatch)
		p.SystemGroupIterator.buffer = make([]ISystem, p.maxPeerBatch)
	}
}

func (p *SystemGroup) systemCount() int {
	return len(p.systems)
}

func (p *SystemGroup) batchCount() int {
	return p.batchTotal
}

func (p *SystemGroup) maxSystemCountPeerBatch() int {
	return p.maxPeerBatch
}

// get all systems
func (p *SystemGroup) all() []ISystem {
	systems := make([]ISystem, len(p.systems))
	for i, n := range p.systems {
		systems[i] = n.val
	}
	return systems
}

// insert system
func (p *SystemGroup) insert(sys ISystem) {
	//get system's required components
	rqs := sys.GetRequirements()
	if len(rqs) == 0 {
		//panic("invalid system")
	}
	//reference count
	for com, _ := range rqs {
		if _, ok := p.ref[com]; ok {
			p.ref[com] += 1
		} else {
			p.ref[com] = 1
		}
	}
	//add system
	node := &Node{
		children: make([]*Node, 0),
		val:      sys,
	}
	p.systems = append(p.systems, node)
	//set unordered
	p.ordered = false
}

// has system
func (p *SystemGroup) has(sys ISystem) bool {
	for _, system := range p.systems {
		if system.val.Type() == sys.Type() {
			return true
		}
	}
	return false
}

// remove system
func (p *SystemGroup) remove(sys ISystem) {
	//get system's required components
	rqs := sys.GetRequirements()
	has := false
	for i, system := range p.systems {
		if system.val.ID() == sys.ID() {
			p.systems = append(p.systems[:i], p.systems[i+1:]...)
			has = true
			break
		}
	}
	if !has {
		return
	}
	//reference count
	for com, _ := range rqs {
		if _, ok := p.ref[com]; ok {
			p.ref[com] -= 1
		} else {
			panic("component ref wrong")
		}
	}
	//set unordered
	p.ordered = false
}

func (p *SystemGroup) iter() *SystemGroupIterator {
	if !p.ordered {
		p.resort()
	}

	if p.maxPeerBatch == 0 {
		return emptySystemGroupIterator
	}

	return &SystemGroupIterator{
		group:   p,
		top:     make([]*Node, p.maxPeerBatch),
		topTemp: make([]*Node, p.maxPeerBatch),
		buffer:  make([]ISystem, p.maxPeerBatch),
	}
}

type SystemGroupIterator struct {
	group   *SystemGroup
	top     []*Node
	topTemp []*Node
	buffer  []ISystem
	topSize int
	size    int
}

func (s *SystemGroupIterator) Begin() []ISystem {
	if s.group == nil {
		return nil
	}
	if !s.group.ordered {
		s.group.resort()
	}
	copy(s.top, s.group.root.children)
	s.topSize = len(s.group.root.children)
	return s.Next()
}

func (s *SystemGroupIterator) Next() []ISystem {
	if s.topSize == 0 {
		s.size = 0
		return nil
	}
	s.topTemp, s.top = s.top, s.topTemp
	tempSize := s.topSize
	s.topSize = 0
	s.size = 0
	for i := 0; i < tempSize; i++ {
		s.buffer[s.size] = s.topTemp[i].val
		s.size++
		for j := 0; j < len(s.topTemp[i].children); j++ {
			s.top[s.topSize+j] = s.topTemp[i].children[j]
		}
		s.topSize += len(s.topTemp[i].children)
	}

	if s.size == 0 {
		return nil
	}
	return s.buffer[:s.size]
}

func (s *SystemGroupIterator) End() bool {
	return s.size == 0
}
