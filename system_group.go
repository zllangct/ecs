package ecs

import (
	"reflect"
	"sort"
	"sync"
)

// system tree node
type Node struct {
	parent   *Node
	children []*Node
	val      ISystem
}

func (p *Node) isFriend(node *Node) bool {
	for com, _ := range p.val.Requirements() {
		for comTarget, _ := range node.val.Requirements() {
			if comTarget == com {
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
		p.children = append(p.children, node)
	}
}

// system group ordered by interrelation
type SystemGroup struct {
	lock    sync.Mutex
	systems []*Node
	ref     map[reflect.Type]int
	top     []*Node
	root    *Node
	ordered bool
	order   Order
}

func NewSystemGroup() *SystemGroup {
	return &SystemGroup{
		lock:    sync.Mutex{},
		systems: make([]*Node, 0),
		ref:     map[reflect.Type]int{},
		ordered: false,
	}
}

func (p *SystemGroup) refCount(rqs map[reflect.Type]struct{}) int {
	ref := 0
	for com, _ := range rqs {
		ref += p.ref[com] - 1
	}
	return ref
}

//initialise system group iterator
func (p *SystemGroup) reset() {
	//need resort
	if !p.ordered {
		sort.Slice(p.systems, func(i, j int) bool {
			return p.refCount(p.systems[i].val.Requirements()) >
				p.refCount(p.systems[j].val.Requirements())
		})
		if p.root == nil {
			p.root = &Node{
				parent:   nil,
				children: []*Node{},
				val:      nil,
			}
		}
		for _, node := range p.systems {
			p.root.attach(node)
		}
		p.ordered = true
	}
	// initialise the iterator
	p.top = make([]*Node, 0)
	if p.root != nil {
		p.top = append(p.top, p.root)
	}
}

//Pop a batch of independent system array
func (p *SystemGroup) next() []ISystem {
	temp := make([]*Node, 0)
	systems := make([]ISystem, 0)
	for _, n := range p.top {
		temp = append(temp, n.children...)
		for _, sys := range n.children {
			systems = append(systems, sys.val)
		}
	}
	p.top = temp
	return systems
}

//get all systems
func (p *SystemGroup) all() []ISystem {
	systems := make([]ISystem, len(p.systems))
	for i, n := range p.systems {
		systems[i] = n.val
	}
	return systems
}

//insert system
func (p *SystemGroup) insert(sys ISystem) {
	//set cluster no ordered
	p.ordered = false
	//get system's required components
	rqs := sys.Requirements()
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
}
