package ecs

import "testing"

func regSizes(its ...ComponentIntType) map[ComponentIntType]int {
	m := map[ComponentIntType]int{}
	for _, it := range its {
		m[it] = 16
	}
	return m
}

func TestArchetypeRegistry_MergeOverlap(t *testing.T) {
	reg := NewArchetypeRegistry()
	// 组1 {1001,1002}，组2 {1002,1003} 重叠 → 并集 {1001,1002,1003}
	if err := reg.Declare([]ComponentIntType{1001, 1002}, regSizes(1001, 1002)); err != nil {
		t.Fatal(err)
	}
	if err := reg.Declare([]ComponentIntType{1002, 1003}, regSizes(1002, 1003)); err != nil {
		t.Fatal(err)
	}
	if err := reg.Finalize(); err != nil {
		t.Fatal(err)
	}
	if len(reg.groups) != 1 {
		t.Fatalf("groups = %d, want 1 merged", len(reg.groups))
	}
	for _, it := range []ComponentIntType{1001, 1002, 1003} {
		if reg.ownerOf[it] == nil {
			t.Fatalf("type %d no owner", it)
		}
	}
}

func TestArchetypeRegistry_MergeTransitive(t *testing.T) {
	reg := NewArchetypeRegistry()
	// {1,2} {3,4} {2,3} → 传递合并 {1,2,3,4}
	reg.Declare([]ComponentIntType{1, 2}, regSizes(1, 2))
	reg.Declare([]ComponentIntType{3, 4}, regSizes(3, 4))
	reg.Declare([]ComponentIntType{2, 3}, regSizes(2, 3))
	if err := reg.Finalize(); err != nil {
		t.Fatal(err)
	}
	if len(reg.groups) != 1 {
		t.Fatalf("groups = %d, want 1 transitive merged", len(reg.groups))
	}
	for it := ComponentIntType(1); it <= 4; it++ {
		if reg.ownerOf[it] == nil {
			t.Fatalf("type %d no owner", it)
		}
	}
}

func TestArchetypeRegistry_DisjointGroups(t *testing.T) {
	reg := NewArchetypeRegistry()
	reg.Declare([]ComponentIntType{1001, 1002}, regSizes(1001, 1002))
	reg.Declare([]ComponentIntType{2001, 2002}, regSizes(2001, 2002))
	if err := reg.Finalize(); err != nil {
		t.Fatal(err)
	}
	if len(reg.groups) != 2 {
		t.Fatalf("groups = %d, want 2", len(reg.groups))
	}
}

func TestArchetypeRegistry_DeclareValidation(t *testing.T) {
	reg := NewArchetypeRegistry()
	if err := reg.Declare([]ComponentIntType{1001}, regSizes(1001)); err == nil {
		t.Fatal("single-type group should fail")
	}
	if err := reg.Declare([]ComponentIntType{1001, 1002}, regSizes(1001)); err == nil {
		t.Fatal("missing size should fail")
	}
	if err := reg.Finalize(); err != nil {
		t.Fatal(err)
	}
	if err := reg.Declare([]ComponentIntType{3, 4}, regSizes(3, 4)); err == nil {
		t.Fatal("declare after finalize should fail")
	}
}
