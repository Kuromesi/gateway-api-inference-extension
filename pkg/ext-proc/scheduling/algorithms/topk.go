package algorithms

import "container/heap"

type TopK[T any] interface {
	TopK(elems []T, k int) []T
}

type HeapTopKImpl[T any] struct {
	cmp    func(a, b T) bool
	sorted []T
}

func NewHeapTopK[T any](cmp func(a, b T) bool) TopK[T] {
	return &HeapTopKImpl[T]{
		cmp: cmp,
	}
}

func (h *HeapTopKImpl[T]) TopK(elems []T, k int) []T {
	if k <= 0 {
		return []T{}
	}

	if k >= len(elems) {
		return elems
	}

	h.sorted = []T{}
	heap.Init(h)
	for _, e := range elems {
		heap.Push(h, e)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return h.sorted
}

func (h *HeapTopKImpl[T]) Len() int           { return len(h.sorted) }
func (h *HeapTopKImpl[T]) Less(i, j int) bool { return h.cmp(h.sorted[i], h.sorted[j]) }
func (h *HeapTopKImpl[T]) Swap(i, j int)      { h.sorted[i], h.sorted[j] = h.sorted[j], h.sorted[i] }

func (h *HeapTopKImpl[T]) Push(x any) {
	h.sorted = append(h.sorted, x.(T))
}

func (h *HeapTopKImpl[T]) Pop() any {
	pop := h.sorted[len(h.sorted)-1]
	h.sorted = h.sorted[:len(h.sorted)-1]
	return pop
}
