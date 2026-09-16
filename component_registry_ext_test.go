package ecs

import (
	"testing"

	rockmem "github.com/zllangct/rockmem/golang"
)

type PositionTestComp struct {
	X, Y, Z float32
	_       [4]byte
}

func (c *PositionTestComp) NewComponentSet() ComponentSet {
	return NewCSet[PositionTestComp]()
}

func (c *PositionTestComp) PacketIdentifier() rockmem.PacketIdentifier {
	return 99999999990001
}

func (c *PositionTestComp) IsNomadic() bool    { return false }
func (c *PositionTestComp) IsDisposable() bool { return false }

func TestComponentDefineInfo_SizeAndProto(t *testing.T) {
	RegisterComponent[PositionTestComp]("test_ext")
	info, ok := GetComponentRegistry().GetInfo(GetIntType[PositionTestComp, *PositionTestComp]())
	if !ok {
		t.Fatal("not registered")
	}
	if info.Size != 16 {
		t.Fatalf("size = %d, want 16", info.Size)
	}
	if info.Proto == nil {
		t.Fatal("proto missing")
	}
	if info.Proto().IsDisposable() || info.Proto().IsNomadic() {
		t.Fatal("proto flags wrong")
	}
}
