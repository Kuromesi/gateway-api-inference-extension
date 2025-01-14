package algorithms

import (
	"reflect"
	"testing"
)

func TestMaxTopK(t *testing.T) {
	tests := []struct {
		elems []int
		k     int
		want  []int
	}{
		{[]int{1, 2, 3, 4, 5}, 3, []int{3, 4, 5}},
		{[]int{5, 4, 3, 2, 1}, 2, []int{4, 5}},
		{[]int{1, 3, 5, 7, 9}, 0, []int{}},
		{[]int{}, 3, []int{}},
		{[]int{10}, 1, []int{10}},
		{[]int{1, 2, 3}, 5, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		h := &HeapTopKImpl[int]{cmp: func(a, b int) bool { return a < b }}
		got := h.TopK(tt.elems, tt.k)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopK(%v, %d) = %v, want %v", tt.elems, tt.k, got, tt.want)
		}
	}
}

func TestMinTopK(t *testing.T) {
	tests := []struct {
		elems []int
		k     int
		want  []int
	}{
		{[]int{1, 2, 3, 4, 5}, 3, []int{3, 1, 2}},
		{[]int{5, 4, 3, 2, 1}, 2, []int{2, 1}},
		{[]int{1, 2, 3}, 5, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		h := &HeapTopKImpl[int]{cmp: func(a, b int) bool { return a > b }}
		got := h.TopK(tt.elems, tt.k)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopK(%v, %d) = %v, want %v", tt.elems, tt.k, got, tt.want)
		}
	}
}
