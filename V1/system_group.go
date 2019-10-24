package main

import (
	"reflect"
	"sort"
	"sync"
)

type Node struct {
	parent *Node
	children []*Node
	val ISystem
}

func (p *Node)isFriend(node *Node) bool {
	for _, com := range p.val.GetRequirements() {
		for _, comTarget := range node.val.GetRequirements() {
			if comTarget.String() == com.String() {
				return true
			}
		}
	}
	return false
}

func (p *Node)attach(node *Node) {
	isAttached := false
	for i := 0; i<len(p.children);i++ {
		if p.children[i].isFriend(node) {
			p.children[i].attach(node)
			if len(p.children[i].children) > len(p.children) {
				temp := p.children[i]
				p.children = append(p.children[:i],p.children[i+1:]...)
				temp.children = append(temp.children, p)
				isAttached = true
				break
			}
		}
	}
	if !isAttached {
		p.children = append(p.children, node)
	}
}

type Group struct {
	val *Node
	ref []string
}

type SystemGroup struct {
	lock sync.Mutex
	systems []*Node
	ref map[string]int
	top []*Node
	root *Node
	ordered bool
	order int
}

func NewSystemGroup() *SystemGroup {
	return &SystemGroup{
		lock: sync.Mutex{},
		systems: make([]*Node,0),
		ref: map[string]int{},
		ordered: false,
	}
}

func (p *SystemGroup)refCount(rqs []reflect.Type) int {
	ref:=0
	for _, com := range rqs {
		ref += p.ref[com.String()]-1
	}
	return ref
}

func (p *SystemGroup)iterInit()  {
	//need resort
	if !p.ordered{
		sort.Slice(p.systems, func(i, j int) bool {
			 return p.refCount(p.systems[i].val.GetRequirements()) >
			 	p.refCount(p.systems[j].val.GetRequirements())
		})

		for _, node := range p.systems {
			if p.root == nil {
				p.root = &Node{
					parent:   nil,
					children: []*Node{},
					val:      nil,
				}
			}else{
				p.root.attach(node)
			}
		}
	}
	// initialise the iterator
	p.top = []*Node{p.root}
}

func (p *SystemGroup)pop()[]ISystem {
	temp := make([]*Node,0)
	systems := make([]ISystem,0)
	for _, n := range p.top {
		temp = append(temp, n.children...)
		for _, sys := range n.children {
			systems = append(systems, sys.val)
		}
	}
	p.top = temp
	return systems
}

func (p *SystemGroup) Insert(sys ISystem) {
	//set cluster no ordered
	p.ordered = false
	//get system's required components
	rqs := sys.GetRequirements()
	if len(rqs) == 0 {
		panic("invalid system")
	}
	//reference count
	for _, com := range rqs {
		if _,ok:=p.ref[com.String()];ok{
			p.ref[com.String()] += 1
		}else{
			p.ref[com.String()] = 1
		}
	}
	//add system
	node := &Node{
		children:make([]*Node,0),
		val:sys,
	}
	p.systems = append(p.systems, node)
}
