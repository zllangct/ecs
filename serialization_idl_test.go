package ecs

import (
	"testing"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// 含 string 字段的非 POD 测试组件，手写最小 rockmem.Message 实现
type stringyComp struct {
	Name string
	Seq  int32
}

func (d *stringyComp) NewComponentSet() ComponentSet {
	return NewCSet[stringyComp]()
}

func (d *stringyComp) PacketIdentifier() rockmem.PacketIdentifier {
	return 65527
}

func (d *stringyComp) IsNomadic() bool {
	return false
}

func (d *stringyComp) IsDisposable() bool {
	return false
}

func (d *stringyComp) WriteAsRoot(w rockmem.Writer) (uint, error) {
	off, err := w.Alloc(uint(8 + len(d.Name)))
	if err != nil {
		return 0, err
	}
	w.Write4At(off, uint32(len(d.Name)))
	w.Write4At(off+4, uint32(d.Seq))
	w.WriteAt(off+8, []byte(d.Name))
	return off, nil
}

func (d *stringyComp) ReadAsRoot(r *rockmem.Reader) {
	l := *(*uint32)(r.Pointer)
	d.Seq = *(*int32)(unsafe.Add(r.Pointer, 4))
	d.Name = string(unsafe.Slice((*byte)(unsafe.Add(r.Pointer, 8)), l))
}

// S1: 非 POD 组件序列化必须逐元素走 IDL，恢复数据不得与原组件共享指针
func TestCSetSerialization_IDLPathNoSharedPointer(t *testing.T) {
	c := NewCSet[stringyComp]()
	e1 := Entity(1 << 0) // index=1
	orig := &stringyComp{Name: "hello-world", Seq: 42}
	c.Add(e1, orig)

	data := c.Marshal()
	restored := NewCSetFromData[stringyComp](data)

	rc := restored.Get(e1)
	if rc == nil {
		t.Fatal("component not restored")
	}
	rs := rc.(*stringyComp)
	if rs.Name != orig.Name || rs.Seq != orig.Seq {
		t.Fatalf("data mismatch: %+v vs %+v", rs, orig)
	}
	// raw 内存直拷会共享 string 底层指针（悬空风险），IDL 路径必须独立
	if unsafe.StringData(rs.Name) == unsafe.StringData(orig.Name) {
		t.Fatal("restored string shares underlying pointer with original (raw copy path)")
	}
}

// S1: world 级快照/恢复整链路（非 POD 组件）
func TestWorldSnapshotRestore_NonPOD(t *testing.T) {
	RegisterComponent[stringyComp]("test")

	w := NewWorld()
	e := w.NewEntity(WithComponents(&stringyComp{Name: "snapshot", Seq: 7}))
	w.Update()

	data := w.Marshal()
	w2 := NewWorldFromData(data)

	ww2 := w2
	info, ok := ww2.getEntityInfo(e)
	if !ok {
		t.Fatal("entity not restored")
	}
	if !info.compound.Exist(GetIntTypeByComp(&stringyComp{})) {
		t.Fatal("compound not restored")
	}
	set, ok := ww2.getComponentSet(GetIntTypeByComp(&stringyComp{}))
	if !ok {
		t.Fatal("component set not restored")
	}
	got := set.Get(e).(*stringyComp)
	if got.Name != "snapshot" || got.Seq != 7 {
		t.Fatalf("component data mismatch: %+v", got)
	}
}

// S1: MarshalTo/UnmarshalFrom 经 rockmem writer 的完整链路
func TestCSetSerialization_IDLPathWithRockmem(t *testing.T) {
	c := NewCSet[stringyComp]()
	for i := 0; i < 3; i++ {
		e := Entity(i + 1)
		c.Add(e, &stringyComp{Name: string(rune('a' + i)), Seq: int32(i)})
	}
	writer := rockmem.NewWriter()
	if _, err := c.MarshalTo(writer); err != nil {
		t.Fatal(err)
	}
	restored := NewCSet[stringyComp]()
	restored.UnmarshalFrom(rockmem.NewReader(writer.Bytes()))
	if restored.Len() != 3 {
		t.Fatalf("want 3, got %d", restored.Len())
	}
	for i := 0; i < 3; i++ {
		e := Entity(i + 1)
		got := restored.Get(e).(*stringyComp)
		if got.Name != string(rune('a'+i)) || got.Seq != int32(i) {
			t.Errorf("entity %d mismatch: %+v", i, got)
		}
	}
}
