package data_structures

import (
	"cmp"
)

// Heap is a generic heap implementation that can be used as min or max heap
type Heap[T any] struct {
	data    []T
	less    func(a, b T) bool
	size    int
}

// NewMinHeap creates a new min heap
func NewMinHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: func(a, b T) bool { return a < b },
		size: 0,
	}
}

// NewMaxHeap creates a new max heap
func NewMaxHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: func(a, b T) bool { return a > b },
		size: 0,
	}
}

// NewHeapWithComparator creates a heap with custom comparator
func NewHeapWithComparator[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: less,
		size: 0,
	}
}

// Push adds an element to the heap
func (h *Heap[T]) Push(val T) {
	h.data = append(h.data, val)
	h.size++
	h.heapifyUp(h.size - 1)
}

// Pop removes and returns the top element
func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	if h.size == 0 {
		return zero, false
	}

	top := h.data[0]
	h.data[0] = h.data[h.size-1]
	h.data = h.data[:h.size-1]
	h.size--

	if h.size > 0 {
		h.heapifyDown(0)
	}

	return top, true
}

// Peek returns the top element without removing it
func (h *Heap[T]) Peek() (T, bool) {
	var zero T
	if h.size == 0 {
		return zero, false
	}
	return h.data[0], true
}

// Size returns the number of elements in the heap
func (h *Heap[T]) Size() int {
	return h.size
}

// IsEmpty returns true if the heap is empty
func (h *Heap[T]) IsEmpty() bool {
	return h.size == 0
}

// heapifyUp maintains heap property by moving element up
func (h *Heap[T]) heapifyUp(idx int) {
	for idx > 0 {
		parent := (idx - 1) / 2
		if h.less(h.data[idx], h.data[parent]) {
			h.data[idx], h.data[parent] = h.data[parent], h.data[idx]
			idx = parent
		} else {
			break
		}
	}
}

// heapifyDown maintains heap property by moving element down
func (h *Heap[T]) heapifyDown(idx int) {
	for {
		left := 2*idx + 1
		right := 2*idx + 2
		smallest := idx

		if left < h.size && h.less(h.data[left], h.data[smallest]) {
			smallest = left
		}

		if right < h.size && h.less(h.data[right], h.data[smallest]) {
			smallest = right
		}

		if smallest != idx {
			h.data[idx], h.data[smallest] = h.data[smallest], h.data[idx]
			idx = smallest
		} else {
			break
		}
	}
}

// BuildHeap creates a heap from a slice in O(n) time
func BuildHeap[T any](data []T, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{
		data: make([]T, len(data)),
		less: less,
		size: len(data),
	}
	copy(h.data, data)

	// Start from last non-leaf node and heapify down
	for i := h.size/2 - 1; i >= 0; i-- {
		h.heapifyDown(i)
	}

	return h
}

// ToSlice returns all elements as a slice (not in sorted order)
func (h *Heap[T]) ToSlice() []T {
	result := make([]T, h.size)
	copy(result, h.data[:h.size])
	return result
}

// Clear removes all elements from the heap
func (h *Heap[T]) Clear() {
	h.data = h.data[:0]
	h.size = 0
}