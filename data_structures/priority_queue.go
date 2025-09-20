package data_structures

import (
	"cmp"
)

// PriorityQueue item with value and priority
type PQItem[T any] struct {
	Value    T
	Priority int
}

// PriorityQueue is a priority queue implementation using heap
type PriorityQueue[T any] struct {
	heap *Heap[PQItem[T]]
}

// NewPriorityQueue creates a new priority queue (min priority first)
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: NewHeapWithComparator(func(a, b PQItem[T]) bool {
			return a.Priority < b.Priority
		}),
	}
}

// NewMaxPriorityQueue creates a new priority queue (max priority first)
func NewMaxPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: NewHeapWithComparator(func(a, b PQItem[T]) bool {
			return a.Priority > b.Priority
		}),
	}
}

// Push adds an item with priority
func (pq *PriorityQueue[T]) Push(value T, priority int) {
	pq.heap.Push(PQItem[T]{Value: value, Priority: priority})
}

// Pop removes and returns the highest priority item
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	var zero T
	item, ok := pq.heap.Pop()
	if !ok {
		return zero, false
	}
	return item.Value, true
}

// Peek returns the highest priority item without removing
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	var zero T
	item, ok := pq.heap.Peek()
	if !ok {
		return zero, false
	}
	return item.Value, true
}

// Size returns the number of items
func (pq *PriorityQueue[T]) Size() int {
	return pq.heap.Size()
}

// IsEmpty returns true if queue is empty
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.heap.IsEmpty()
}

// Clear removes all items
func (pq *PriorityQueue[T]) Clear() {
	pq.heap.Clear()
}

// IndexedHeap for efficient updates/deletes by maintaining index map
type IndexedHeap[T comparable] struct {
	data     []T
	indexMap map[T]int
	less     func(a, b T) bool
	size     int
}

// NewIndexedMinHeap creates an indexed min heap
func NewIndexedMinHeap[T cmp.Ordered]() *IndexedHeap[T] {
	return &IndexedHeap[T]{
		data:     make([]T, 0),
		indexMap: make(map[T]int),
		less:     func(a, b T) bool { return a < b },
		size:     0,
	}
}

// Push adds element to indexed heap
func (h *IndexedHeap[T]) Push(val T) {
	if _, exists := h.indexMap[val]; exists {
		return // Already exists
	}

	h.data = append(h.data, val)
	h.indexMap[val] = h.size
	h.size++
	h.heapifyUp(h.size - 1)
}

// Update changes the value and maintains heap property
func (h *IndexedHeap[T]) Update(oldVal, newVal T) bool {
	idx, exists := h.indexMap[oldVal]
	if !exists {
		return false
	}

	delete(h.indexMap, oldVal)
	h.data[idx] = newVal
	h.indexMap[newVal] = idx

	// Try both directions as we don't know if value increased or decreased
	h.heapifyUp(idx)
	h.heapifyDown(idx)
	return true
}

// Delete removes a specific value from heap
func (h *IndexedHeap[T]) Delete(val T) bool {
	idx, exists := h.indexMap[val]
	if !exists {
		return false
	}

	delete(h.indexMap, val)

	if idx == h.size-1 {
		h.data = h.data[:h.size-1]
		h.size--
		return true
	}

	// Move last element to deleted position
	h.data[idx] = h.data[h.size-1]
	h.indexMap[h.data[idx]] = idx
	h.data = h.data[:h.size-1]
	h.size--

	if h.size > 0 {
		h.heapifyUp(idx)
		h.heapifyDown(idx)
	}

	return true
}

// Contains checks if value exists
func (h *IndexedHeap[T]) Contains(val T) bool {
	_, exists := h.indexMap[val]
	return exists
}

// Pop removes and returns the top element
func (h *IndexedHeap[T]) Pop() (T, bool) {
	var zero T
	if h.size == 0 {
		return zero, false
	}

	top := h.data[0]
	delete(h.indexMap, top)

	if h.size == 1 {
		h.data = h.data[:0]
		h.size = 0
		return top, true
	}

	h.data[0] = h.data[h.size-1]
	h.indexMap[h.data[0]] = 0
	h.data = h.data[:h.size-1]
	h.size--
	h.heapifyDown(0)

	return top, true
}

// heapifyUp for indexed heap
func (h *IndexedHeap[T]) heapifyUp(idx int) {
	for idx > 0 {
		parent := (idx - 1) / 2
		if h.less(h.data[idx], h.data[parent]) {
			h.swap(idx, parent)
			idx = parent
		} else {
			break
		}
	}
}

// heapifyDown for indexed heap
func (h *IndexedHeap[T]) heapifyDown(idx int) {
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
			h.swap(idx, smallest)
			idx = smallest
		} else {
			break
		}
	}
}

// swap exchanges elements and updates index map
func (h *IndexedHeap[T]) swap(i, j int) {
	h.indexMap[h.data[i]], h.indexMap[h.data[j]] = j, i
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Size returns the number of elements
func (h *IndexedHeap[T]) Size() int {
	return h.size
}

// IsEmpty returns true if heap is empty
func (h *IndexedHeap[T]) IsEmpty() bool {
	return h.size == 0
}