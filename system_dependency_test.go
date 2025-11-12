package ecs

import (
	"testing"
)

func TestDep(t *testing.T) {
	dep1 := Dep[dummyComponent](ReadOnly)
	if dep1.intType() != GetIntType[dummyComponent]() {
		t.Errorf("int type not equal")
	}
	if dep1.readonly() != true {
		t.Errorf("readonly not equal")
	}

	dep2 := Dep[dummyComponent](ReadWrite)
	if dep2.intType() != GetIntType[dummyComponent]() {
		t.Errorf("int type not equal")
	}
	if dep2.readonly() != false {
		t.Errorf("readonly not equal")
	}
}
