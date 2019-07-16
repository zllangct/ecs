package main

type NameComponent struct {
	ComponentBase
	Name string
}

type CountComponent struct {
	ComponentBase
	Count int
}

type AddressComponent struct {
	ComponentBase
	Address string
}

type NameSystem struct {
	SystemBase
}

func (this *NameSystem) Init() {

}

func (*NameSystem) ComponentRequire() []IComponent {
	return []IComponent{&NameComponent{}}
}

func (*NameSystem) Update() {
	panic("implement me")
}

type OtherSystem struct {
	SystemBase
}

func (*OtherSystem) Init() {

}

func (*OtherSystem) ComponentRequire() []IComponent {
	return []IComponent{&AddressComponent{},&CountComponent{}}
}

func (this *OtherSystem) Update() {
	//for _, com := range this.components {
	//	otherc:=()com
	//}
}

func main() {
	//runtime := Runtime{
	//	UpdateInterval:500,
	//}
	//
	//runtime.Filter(&NameSystem{})
	//runtime.Filter(&OtherSystem{})
	//
	//runtime.Run()

	m:=map[interface{}]int{
		struct {
			Name string
		}{}:2,
		struct {
			Count int
		}{}:3,
	}

	println(m[struct {
		Count int
	}{}])
	s,ok:=m["sss"]
	println(ok,s)
}
