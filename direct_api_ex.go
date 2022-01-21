package ecs

func AddRequireComponent2[T1 ComponentObject, T2 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2]())
}

func AddRequireComponent3[T1 ComponentObject, T2 ComponentObject, T3 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3]())
}

func AddRequireComponent4[T1 ComponentObject, T2 ComponentObject, T3 ComponentObject, T4 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4]())
}

func AddRequireComponent5[T1 ComponentObject, T2 ComponentObject, T3 ComponentObject, T4 ComponentObject, T5 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5]())
}

func AddRequireComponent6[T1 ComponentObject, T2 ComponentObject, T3 ComponentObject, T4 ComponentObject, T5 ComponentObject, T6 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6]())
}

func AddRequireComponent7[T1 ComponentObject, T2 ComponentObject, T3 ComponentObject, T4 ComponentObject, T5 ComponentObject, T6 ComponentObject, T7 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6](), TypeOf[T7]())
}

func AddRequireComponent8[T1 ComponentObject, T2 ComponentObject, T3 ComponentObject, T4 ComponentObject, T5 ComponentObject, T6 ComponentObject, T7 ComponentObject, T8 ComponentObject](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T1](), TypeOf[T2](), TypeOf[T3](), TypeOf[T4](), TypeOf[T5](), TypeOf[T6](), TypeOf[T7](), TypeOf[T8]())
}
