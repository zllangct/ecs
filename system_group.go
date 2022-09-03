package ecs

import (
	"reflect"
	"sort"
	"sync"
)

var emptyGroup []ISystem

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
		ordered: true,
	}
}

func (p *SystemGroup) refCount(rqs map[reflect.Type]IRequirement) int {
	ref := 0
	for com, _ := range rqs {
		ref += p.ref[com] - 1
	}
	return ref
}

// initialise system group iterator
func (p *SystemGroup) reset() {
	//need resort
	if !p.ordered {
		sort.Slice(p.systems, func(i, j int) bool {
			return p.refCount(p.systems[i].val.GetRequirements()) >
				p.refCount(p.systems[j].val.GetRequirements())
		})
		if p.root == nil {
			p.root = &Node{
				parent:   nil,
				children: []*Node{},
				val:      nil,
			}
		}

		p.root.children = []*Node{}
		for _, node := range p.systems {
			node.children = []*Node{}
			p.root.attach(node)
		}
		p.ordered = true
	}

	if len(p.systems) == 0 {
		return
	}

	// initialise the iterator
	p.top = p.root.children
}

// Pop a batch of independent system array
func (p *SystemGroup) next() []ISystem {
	if p.top == nil {
		return emptyGroup
	}
	systems := make([]ISystem, 0)
	temp := p.top
	p.top = make([]*Node, 0)
	for _, s := range temp {
		systems = append(systems, s.val)
		for _, n := range s.children {
			p.top = append(p.top, n)
		}
	}
	return systems
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
	//set cluster no ordered
	p.ordered = false
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
	if len(rqs) == 0 {
		//panic("invalid system")
	}
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
	//set cluster no ordered
	p.ordered = false
	//reference count
	for com, _ := range rqs {
		if _, ok := p.ref[com]; ok {
			p.ref[com] -= 1
		} else {
			println("component ref wrong")
		}
	}

	p.reset()
}
