package ecs

//func AttachTo[C any, T ITComponent[C]](e *Entity, com ... T) {
//	if len(com) == 0 {
//		ins := new(C)
//		e.AddComponent(ins)
//	} else {
//		e.AddComponent(com...)
//	}
//}
//
//func GetComponentFrom[C any, T ITComponent[C]](e *Entity) *C{
//	return (*C)(e.GetComponent(reflect.TypeOf(*new(C))))
//}
