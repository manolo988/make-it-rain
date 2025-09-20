package data_structures

import (
	"cmp"
)

// Stack is a generic stack implementation
type Stack[T any] struct {
	data []T
	size int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0),
		size: 0,
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(val T) {
	s.data = append(s.data, val)
	s.size++
}

// Pop removes and returns the top element
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if s.size == 0 {
		return zero, false
	}

	val := s.data[s.size-1]
	s.data = s.data[:s.size-1]
	s.size--
	return val, true
}

// Peek returns the top element without removing it
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if s.size == 0 {
		return zero, false
	}
	return s.data[s.size-1], true
}

// IsEmpty returns true if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return s.size == 0
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return s.size
}

// Clear removes all elements from the stack
func (s *Stack[T]) Clear() {
	s.data = s.data[:0]
	s.size = 0
}

// ToSlice returns all elements as a slice (bottom to top)
func (s *Stack[T]) ToSlice() []T {
	result := make([]T, s.size)
	copy(result, s.data)
	return result
}

// MinStack maintains minimum element in O(1) time
type MinStack[T cmp.Ordered] struct {
	data     []T
	minStack []T
	size     int
}

// NewMinStack creates a new MinStack
func NewMinStack[T cmp.Ordered]() *MinStack[T] {
	return &MinStack[T]{
		data:     make([]T, 0),
		minStack: make([]T, 0),
		size:     0,
	}
}

// Push adds element and updates minimum
func (ms *MinStack[T]) Push(val T) {
	ms.data = append(ms.data, val)

	if len(ms.minStack) == 0 || val <= ms.minStack[len(ms.minStack)-1] {
		ms.minStack = append(ms.minStack, val)
	}
	ms.size++
}

// Pop removes top element and updates minimum
func (ms *MinStack[T]) Pop() (T, bool) {
	var zero T
	if ms.size == 0 {
		return zero, false
	}

	val := ms.data[ms.size-1]
	ms.data = ms.data[:ms.size-1]

	if len(ms.minStack) > 0 && val == ms.minStack[len(ms.minStack)-1] {
		ms.minStack = ms.minStack[:len(ms.minStack)-1]
	}

	ms.size--
	return val, true
}

// Peek returns top element
func (ms *MinStack[T]) Peek() (T, bool) {
	var zero T
	if ms.size == 0 {
		return zero, false
	}
	return ms.data[ms.size-1], true
}

// GetMin returns minimum element in O(1)
func (ms *MinStack[T]) GetMin() (T, bool) {
	var zero T
	if len(ms.minStack) == 0 {
		return zero, false
	}
	return ms.minStack[len(ms.minStack)-1], true
}

// IsEmpty returns true if stack is empty
func (ms *MinStack[T]) IsEmpty() bool {
	return ms.size == 0
}

// Size returns number of elements
func (ms *MinStack[T]) Size() int {
	return ms.size
}

// MaxStack maintains maximum element in O(1) time
type MaxStack[T cmp.Ordered] struct {
	data     []T
	maxStack []T
	size     int
}

// NewMaxStack creates a new MaxStack
func NewMaxStack[T cmp.Ordered]() *MaxStack[T] {
	return &MaxStack[T]{
		data:     make([]T, 0),
		maxStack: make([]T, 0),
		size:     0,
	}
}

// Push adds element and updates maximum
func (ms *MaxStack[T]) Push(val T) {
	ms.data = append(ms.data, val)

	if len(ms.maxStack) == 0 || val >= ms.maxStack[len(ms.maxStack)-1] {
		ms.maxStack = append(ms.maxStack, val)
	}
	ms.size++
}

// Pop removes top element and updates maximum
func (ms *MaxStack[T]) Pop() (T, bool) {
	var zero T
	if ms.size == 0 {
		return zero, false
	}

	val := ms.data[ms.size-1]
	ms.data = ms.data[:ms.size-1]

	if len(ms.maxStack) > 0 && val == ms.maxStack[len(ms.maxStack)-1] {
		ms.maxStack = ms.maxStack[:len(ms.maxStack)-1]
	}

	ms.size--
	return val, true
}

// GetMax returns maximum element in O(1)
func (ms *MaxStack[T]) GetMax() (T, bool) {
	var zero T
	if len(ms.maxStack) == 0 {
		return zero, false
	}
	return ms.maxStack[len(ms.maxStack)-1], true
}

// MonotonicStack maintains elements in monotonic order
type MonotonicStack[T cmp.Ordered] struct {
	data       []T
	increasing bool // true for increasing, false for decreasing
}

// NewMonotonicStack creates a monotonic stack
func NewMonotonicStack[T cmp.Ordered](increasing bool) *MonotonicStack[T] {
	return &MonotonicStack[T]{
		data:       make([]T, 0),
		increasing: increasing,
	}
}

// Push maintains monotonic property
func (ms *MonotonicStack[T]) Push(val T) []T {
	popped := make([]T, 0)

	if ms.increasing {
		// Pop all elements greater than val
		for len(ms.data) > 0 && ms.data[len(ms.data)-1] > val {
			popped = append(popped, ms.data[len(ms.data)-1])
			ms.data = ms.data[:len(ms.data)-1]
		}
	} else {
		// Pop all elements smaller than val
		for len(ms.data) > 0 && ms.data[len(ms.data)-1] < val {
			popped = append(popped, ms.data[len(ms.data)-1])
			ms.data = ms.data[:len(ms.data)-1]
		}
	}

	ms.data = append(ms.data, val)
	return popped
}

// Pop removes and returns top element
func (ms *MonotonicStack[T]) Pop() (T, bool) {
	var zero T
	if len(ms.data) == 0 {
		return zero, false
	}

	val := ms.data[len(ms.data)-1]
	ms.data = ms.data[:len(ms.data)-1]
	return val, true
}

// Peek returns top element without removing
func (ms *MonotonicStack[T]) Peek() (T, bool) {
	var zero T
	if len(ms.data) == 0 {
		return zero, false
	}
	return ms.data[len(ms.data)-1], true
}

// IsEmpty returns true if stack is empty
func (ms *MonotonicStack[T]) IsEmpty() bool {
	return len(ms.data) == 0
}

// Size returns number of elements
func (ms *MonotonicStack[T]) Size() int {
	return len(ms.data)
}

// ToSlice returns all elements (bottom to top)
func (ms *MonotonicStack[T]) ToSlice() []T {
	result := make([]T, len(ms.data))
	copy(result, ms.data)
	return result
}