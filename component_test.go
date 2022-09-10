package ecs

import "testing"

func TestComponent_isValidComponentType(t *testing.T) {
	type C1 struct {
		Component[C1]
		Field1 int
		Field2 struct {
			Field1 int
		}
	}

	type C2 struct {
		Component[C2]
		Field1 string
	}

	type C3 struct {
		Component[C3]
		Field1 *int
	}

	type C4 struct {
		Component[C4]
		Field1 int
		Field2 struct {
			Field1 *int
		}
	}

	type C5 struct {
		Component[C5]
		Field1 int
		Field2 struct {
			Field1 struct {
				Field1 int
			}
			Field2 uint32
		}
		Field3 FixedString[Fixed5]
	}

	tests := []struct {
		name string
		c    IComponent
		want bool
	}{
		{
			name: "Test1",
			c:    &C1{},
			want: true,
		},
		{
			name: "Test2",
			c:    &C2{},
			want: false,
		},
		{
			name: "Test3",
			c:    &C3{},
			want: false,
		},
		{
			name: "Test4",
			c:    &C4{},
			want: false,
		},
		{
			name: "Test5",
			c:    &C5{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.isValidComponentType(); got != tt.want {
				t.Errorf("isValidComponentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
