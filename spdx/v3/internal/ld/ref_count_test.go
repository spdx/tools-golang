package ld

import (
	"testing"
)

func Test_refCount(t *testing.T) {
	type O1 struct {
		Name string
	}

	type O2 struct {
		Name string
		O1s  []*O1
	}

	o1 := &O1{"o1"}
	o2 := &O1{"o2"}
	o3 := &O1{"o3"}
	o21 := &O2{"o21", []*O1{o1, o1, o2, o3}}
	o22 := []*O2{
		{"o22-1", []*O1{o1, o1, o1, o1, o2, o3}},
		{"o22-2", []*O1{o1, o1, o1, o1, o2, o3}},
		{"o22-3", []*O1{o1, o1, o1, o1, o2, o3}},
	}

	type O3 struct {
		Name string
		Ref  []*O3
	}
	o31 := &O3{"o31", nil}
	o32 := &O3{"o32", []*O3{o31}}
	o33 := &O3{"o33", []*O3{o32}}
	o31.Ref = []*O3{o33}
	o34 := &O3{"o34", []*O3{o31, o32}}
	o35 := &O3{"o35", []*O3{o31, o32, o31, o32}}

	type O4 struct {
		Name string
		Ref  any
	}
	o41 := &O4{"o41", nil}
	o42 := &O4{"o42", o41}

	tests := []struct {
		name     string
		checkObj any
		checkIn  any
		expected int
	}{
		{
			name:     "none",
			checkObj: o33,
			checkIn:  o21,
			expected: 0,
		},
		{
			name:     "interface",
			checkObj: o41,
			checkIn:  o42,
			expected: 1,
		},
		{
			name:     "single",
			checkObj: o3,
			checkIn:  o21,
			expected: 1,
		},
		{
			name:     "multiple",
			checkObj: o1,
			checkIn:  o21,
			expected: 2,
		},

		{
			name:     "multiple 2",
			checkObj: o1,
			checkIn:  o22,
			expected: 12,
		},
		{
			name:     "circular 1",
			checkObj: o31,
			checkIn:  o31,
			expected: 2, // this returns 2 because it needs to find a circular reference to itself
		},
		{
			name:     "circular 2",
			checkObj: o32,
			checkIn:  o31,
			expected: 1,
		},
		{
			name:     "circular 3",
			checkObj: o33,
			checkIn:  o31,
			expected: 1,
		},
		{
			name:     "circular multiple",
			checkObj: o32,
			checkIn:  o34,
			expected: 2,
		},
		{
			name:     "circular multiple 2",
			checkObj: o32,
			checkIn:  o35,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnt := RefCount(tt.checkObj, tt.checkIn)
			if cnt != tt.expected {
				t.Errorf("wrong reference count: %v != %v", tt.expected, cnt)
			}
		})
	}
}
