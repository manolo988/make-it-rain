package data_structures

import (
	"sync"
)

// Queue is a generic FIFO queue implementation
type Queue[T any] struct {
	data []T
	size int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		data: make([]T, 0),
		size: 0,
	}
}

// Enqueue adds an element to the back of the queue
func (q *Queue[T]) Enqueue(val T) {
	q.data = append(q.data, val)
	q.size++
}

// Dequeue removes and returns the front element
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}

	val := q.data[0]
	q.data = q.data[1:]
	q.size--
	return val, true
}

// Front returns the front element without removing it
func (q *Queue[T]) Front() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.data[0], true
}

// Back returns the back element without removing it
func (q *Queue[T]) Back() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.data[q.size-1], true
}

// IsEmpty returns true if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return q.size
}

// Clear removes all elements from the queue
func (q *Queue[T]) Clear() {
	q.data = q.data[:0]
	q.size = 0
}

// ToSlice returns all elements as a slice (front to back)
func (q *Queue[T]) ToSlice() []T {
	result := make([]T, q.size)
	copy(result, q.data)
	return result
}

// CircularQueue is a fixed-size queue using circular buffer
type CircularQueue[T any] struct {
	data     []T
	front    int
	rear     int
	size     int
	capacity int
}

// NewCircularQueue creates a new circular queue with fixed capacity
func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	return &CircularQueue[T]{
		data:     make([]T, capacity),
		front:    0,
		rear:     0,
		size:     0,
		capacity: capacity,
	}
}

// Enqueue adds element to the queue
func (cq *CircularQueue[T]) Enqueue(val T) bool {
	if cq.size == cq.capacity {
		return false // Queue is full
	}

	cq.data[cq.rear] = val
	cq.rear = (cq.rear + 1) % cq.capacity
	cq.size++
	return true
}

// Dequeue removes and returns front element
func (cq *CircularQueue[T]) Dequeue() (T, bool) {
	var zero T
	if cq.size == 0 {
		return zero, false
	}

	val := cq.data[cq.front]
	cq.front = (cq.front + 1) % cq.capacity
	cq.size--
	return val, true
}

// Front returns the front element without removing
func (cq *CircularQueue[T]) Front() (T, bool) {
	var zero T
	if cq.size == 0 {
		return zero, false
	}
	return cq.data[cq.front], true
}

// IsEmpty returns true if queue is empty
func (cq *CircularQueue[T]) IsEmpty() bool {
	return cq.size == 0
}

// IsFull returns true if queue is full
func (cq *CircularQueue[T]) IsFull() bool {
	return cq.size == cq.capacity
}

// Size returns current number of elements
func (cq *CircularQueue[T]) Size() int {
	return cq.size
}

// Capacity returns the maximum capacity
func (cq *CircularQueue[T]) Capacity() int {
	return cq.capacity
}

// Deque is a double-ended queue
type Deque[T any] struct {
	data []T
	size int
}

// NewDeque creates a new double-ended queue
func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{
		data: make([]T, 0),
		size: 0,
	}
}

// PushFront adds element to the front
func (d *Deque[T]) PushFront(val T) {
	d.data = append([]T{val}, d.data...)
	d.size++
}

// PushBack adds element to the back
func (d *Deque[T]) PushBack(val T) {
	d.data = append(d.data, val)
	d.size++
}

// PopFront removes and returns front element
func (d *Deque[T]) PopFront() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}

	val := d.data[0]
	d.data = d.data[1:]
	d.size--
	return val, true
}

// PopBack removes and returns back element
func (d *Deque[T]) PopBack() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}

	val := d.data[d.size-1]
	d.data = d.data[:d.size-1]
	d.size--
	return val, true
}

// Front returns front element without removing
func (d *Deque[T]) Front() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}
	return d.data[0], true
}

// Back returns back element without removing
func (d *Deque[T]) Back() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}
	return d.data[d.size-1], true
}

// IsEmpty returns true if deque is empty
func (d *Deque[T]) IsEmpty() bool {
	return d.size == 0
}

// Size returns number of elements
func (d *Deque[T]) Size() int {
	return d.size
}

// ThreadSafeQueue is a concurrent-safe queue
type ThreadSafeQueue[T any] struct {
	queue *Queue[T]
	mu    sync.RWMutex
}

// NewThreadSafeQueue creates a thread-safe queue
func NewThreadSafeQueue[T any]() *ThreadSafeQueue[T] {
	return &ThreadSafeQueue[T]{
		queue: NewQueue[T](),
	}
}

// Enqueue adds element thread-safely
func (tsq *ThreadSafeQueue[T]) Enqueue(val T) {
	tsq.mu.Lock()
	defer tsq.mu.Unlock()
	tsq.queue.Enqueue(val)
}

// Dequeue removes element thread-safely
func (tsq *ThreadSafeQueue[T]) Dequeue() (T, bool) {
	tsq.mu.Lock()
	defer tsq.mu.Unlock()
	return tsq.queue.Dequeue()
}

// Front returns front element thread-safely
func (tsq *ThreadSafeQueue[T]) Front() (T, bool) {
	tsq.mu.RLock()
	defer tsq.mu.RUnlock()
	return tsq.queue.Front()
}

// Size returns size thread-safely
func (tsq *ThreadSafeQueue[T]) Size() int {
	tsq.mu.RLock()
	defer tsq.mu.RUnlock()
	return tsq.queue.Size()
}

// IsEmpty returns if empty thread-safely
func (tsq *ThreadSafeQueue[T]) IsEmpty() bool {
	tsq.mu.RLock()
	defer tsq.mu.RUnlock()
	return tsq.queue.IsEmpty()
}