package ecs

func AddRequireComponent2[T1 IComponent, T2 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2]())
}

func AddRequireComponent3[T1 IComponent, T2 IComponent, T3 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3]())
}

func AddRequireComponent4[T1 IComponent, T2 IComponent, T3 IComponent, T4 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4]())
}

func AddRequireComponent5[T1 IComponent, T2 IComponent, T3 IComponent, T4 IComponent, T5 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5]())
}

func AddRequireComponent6[T1 IComponent, T2 IComponent, T3 IComponent, T4 IComponent, T5 IComponent, T6 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6]())
}

func AddRequireComponent7[T1 IComponent, T2 IComponent, T3 IComponent, T4 IComponent, T5 IComponent, T6 IComponent, T7 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6](), TypeOf[T7]())
}

func AddRequireComponent8[T1 IComponent, T2 IComponent, T3 IComponent, T4 IComponent, T5 IComponent, T6 IComponent, T7 IComponent, T8 IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6](), TypeOf[T7](), TypeOf[T8]())
}
