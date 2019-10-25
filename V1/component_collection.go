package main

import "reflect"


type componentData  struct {
	data []interface{}
	index map[string]int
}

func (p *componentData) push(com IComponent,id string) {
	p.data = append(p.data, com)
	p.index[id] = len(p.data)-1
}

func (p *componentData) pop(id string) {
	if index,ok := p.index[id];ok{
		length:=len(p.data)
		p.data[index],p.data[length-1] = p.data[length-1], p.data[index]
		if length > 0 {
			p.data = p.data[:length-1]
		}
	}
}

type ComponentCollection map[reflect.Type]*componentData

func (p ComponentCollection) push(com IComponent,id string)  {
	typ := reflect.TypeOf(com)
	if v,ok:=p[typ];ok {
		v.push(com,id)
	}else{
		cd := &componentData{
			data:  []interface{}{},
			index: map[string]int{},
		}
		cd.push(com,id)
		p[typ]=cd
	}
}

func (p ComponentCollection) pop(id string,typ reflect.Type)  {
	if v,ok:=p[typ];ok {
		v.pop(id)
	}
}

func (p ComponentCollection) GetComponents(typ reflect.Type) []interface{} {
	v,ok:= p[typ]
	if ok {
		return v.data
	}
	return []interface{}{}
}

func (p ComponentCollection) GetComponent(typ reflect.Type,id string) interface{} {
	v,ok:= p[typ]
	if ok {
		if c,ok:=v.index[id]; ok {
			return v.data[c]
		}
	}
	return nil
}

func (p ComponentCollection) GetIterator() *ComponentCollectionIter {
	ls := make([]*componentData,0)
	for _, value := range p {
		ls = append(ls, value)
	}
	return newComponentCollectionIter(ls)
}




