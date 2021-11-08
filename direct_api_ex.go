package ecs

func AddRequireComponent2[T1 IComponentTemplate, T2 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2]())
}

func AddRequireComponent3[T1 IComponentTemplate, T2 IComponentTemplate, T3 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3]())
}

func AddRequireComponent4[T1 IComponentTemplate, T2 IComponentTemplate, T3 IComponentTemplate, T4 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4]())
}

func AddRequireComponent5[T1 IComponentTemplate, T2 IComponentTemplate, T3 IComponentTemplate, T4 IComponentTemplate, T5 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5]())
}

func AddRequireComponent6[T1 IComponentTemplate, T2 IComponentTemplate, T3 IComponentTemplate, T4 IComponentTemplate, T5 IComponentTemplate, T6 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6]())
}

func AddRequireComponent7[T1 IComponentTemplate, T2 IComponentTemplate, T3 IComponentTemplate, T4 IComponentTemplate, T5 IComponentTemplate, T6 IComponentTemplate, T7 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6](), TypeOf[T7]())
}

func AddRequireComponent8[T1 IComponentTemplate, T2 IComponentTemplate, T3 IComponentTemplate, T4 IComponentTemplate, T5 IComponentTemplate, T6 IComponentTemplate, T7 IComponentTemplate, T8 IComponentTemplate](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6](), TypeOf[T7](), TypeOf[T8]())
}
