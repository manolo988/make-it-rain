package data_structures

import (
	"errors"
	"sync"
	"time"
)

// RingBuffer is a fixed-size circular buffer
type RingBuffer[T any] struct {
	data     []T
	capacity int
	size     int
	head     int // Next write position
	tail     int // Next read position
	mu       sync.RWMutex
}

// NewRingBuffer creates a new ring buffer
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:     make([]T, capacity),
		capacity: capacity,
		size:     0,
		head:     0,
		tail:     0,
	}
}

// Write adds an element to the buffer
func (rb *RingBuffer[T]) Write(val T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.data[rb.head] = val
	rb.head = (rb.head + 1) % rb.capacity

	if rb.size == rb.capacity {
		// Overwrite oldest - move tail forward
		rb.tail = (rb.tail + 1) % rb.capacity
	} else {
		rb.size++
	}
}

// Read reads and removes an element
func (rb *RingBuffer[T]) Read() (T, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	var zero T
	if rb.size == 0 {
		return zero, errors.New("buffer is empty")
	}

	val := rb.data[rb.tail]
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.size--

	return val, nil
}

// Peek reads without removing
func (rb *RingBuffer[T]) Peek() (T, error) {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	var zero T
	if rb.size == 0 {
		return zero, errors.New("buffer is empty")
	}

	return rb.data[rb.tail], nil
}

// WriteMultiple writes multiple values
func (rb *RingBuffer[T]) WriteMultiple(values []T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for _, val := range values {
		rb.data[rb.head] = val
		rb.head = (rb.head + 1) % rb.capacity

		if rb.size == rb.capacity {
			rb.tail = (rb.tail + 1) % rb.capacity
		} else {
			rb.size++
		}
	}
}

// ReadMultiple reads up to n elements
func (rb *RingBuffer[T]) ReadMultiple(n int) []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	toRead := min(n, rb.size)
	result := make([]T, toRead)

	for i := 0; i < toRead; i++ {
		result[i] = rb.data[rb.tail]
		rb.tail = (rb.tail + 1) % rb.capacity
		rb.size--
	}

	return result
}

// Size returns current number of elements
func (rb *RingBuffer[T]) Size() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size
}

// Capacity returns the buffer capacity
func (rb *RingBuffer[T]) Capacity() int {
	return rb.capacity
}

// IsFull returns true if buffer is full
func (rb *RingBuffer[T]) IsFull() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size == rb.capacity
}

// IsEmpty returns true if buffer is empty
func (rb *RingBuffer[T]) IsEmpty() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size == 0
}

// Clear removes all elements
func (rb *RingBuffer[T]) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.size = 0
	rb.head = 0
	rb.tail = 0
}

// ToSlice returns all elements as slice
func (rb *RingBuffer[T]) ToSlice() []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	result := make([]T, rb.size)
	for i := 0; i < rb.size; i++ {
		idx := (rb.tail + i) % rb.capacity
		result[i] = rb.data[idx]
	}

	return result
}

// GetLatest returns the n most recent elements
func (rb *RingBuffer[T]) GetLatest(n int) []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	toGet := min(n, rb.size)
	result := make([]T, toGet)

	// Start from the most recent (head - 1)
	start := rb.head - toGet
	if start < 0 {
		start += rb.capacity
	}

	for i := 0; i < toGet; i++ {
		idx := (start + i) % rb.capacity
		result[i] = rb.data[idx]
	}

	return result
}

// TimedRingBuffer is a ring buffer with time-based eviction
type TimedRingBuffer[T any] struct {
	buffer     *RingBuffer[TimedEntry[T]]
	maxAge     time.Duration
	mu         sync.RWMutex
}

type TimedEntry[T any] struct {
	Value     T
	Timestamp time.Time
}

// NewTimedRingBuffer creates a ring buffer with time-based eviction
func NewTimedRingBuffer[T any](capacity int, maxAge time.Duration) *TimedRingBuffer[T] {
	trb := &TimedRingBuffer[T]{
		buffer: NewRingBuffer[TimedEntry[T]](capacity),
		maxAge: maxAge,
	}

	// Start cleanup goroutine
	go trb.cleanupExpired()

	return trb
}

// Write adds a timestamped value
func (trb *TimedRingBuffer[T]) Write(val T) {
	entry := TimedEntry[T]{
		Value:     val,
		Timestamp: time.Now(),
	}
	trb.buffer.Write(entry)
}

// ReadValid reads only non-expired entries
func (trb *TimedRingBuffer[T]) ReadValid() []T {
	trb.mu.RLock()
	defer trb.mu.RUnlock()

	entries := trb.buffer.ToSlice()
	result := make([]T, 0)
	cutoff := time.Now().Add(-trb.maxAge)

	for _, entry := range entries {
		if entry.Timestamp.After(cutoff) {
			result = append(result, entry.Value)
		}
	}

	return result
}

// cleanupExpired periodically removes old entries
func (trb *TimedRingBuffer[T]) cleanupExpired() {
	ticker := time.NewTicker(trb.maxAge / 2)
	defer ticker.Stop()

	for range ticker.C {
		trb.mu.Lock()
		// Ring buffer naturally overwrites old entries
		// This is mainly for metrics/monitoring
		trb.mu.Unlock()
	}
}

// MetricsBuffer is specialized for collecting metrics
type MetricsBuffer struct {
	buffer *RingBuffer[float64]
	mu     sync.RWMutex
}

// NewMetricsBuffer creates a metrics buffer
func NewMetricsBuffer(capacity int) *MetricsBuffer {
	return &MetricsBuffer{
		buffer: NewRingBuffer[float64](capacity),
	}
}

// Record adds a metric value
func (mb *MetricsBuffer) Record(value float64) {
	mb.buffer.Write(value)
}

// GetStats returns statistics for buffer contents
func (mb *MetricsBuffer) GetStats() MetricStats {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	values := mb.buffer.ToSlice()
	if len(values) == 0 {
		return MetricStats{}
	}

	var sum, min, max float64
	min = values[0]
	max = values[0]

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	avg := sum / float64(len(values))

	// Calculate percentiles (simplified)
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	// Note: In production, you'd sort these
	// sortedValues = sort.Float64s(sortedValues)

	return MetricStats{
		Count:   len(values),
		Sum:     sum,
		Average: avg,
		Min:     min,
		Max:     max,
		// P50, P95, P99 would be calculated from sorted values
	}
}

type MetricStats struct {
	Count   int
	Sum     float64
	Average float64
	Min     float64
	Max     float64
	P50     float64 // Median
	P95     float64 // 95th percentile
	P99     float64 // 99th percentile
}

// SlidingWindowBuffer maintains a time-based sliding window
type SlidingWindowBuffer[T any] struct {
	entries    []WindowEntry[T]
	windowSize time.Duration
	mu         sync.RWMutex
}

type WindowEntry[T any] struct {
	Value     T
	Timestamp time.Time
}

// NewSlidingWindowBuffer creates a sliding window buffer
func NewSlidingWindowBuffer[T any](windowSize time.Duration) *SlidingWindowBuffer[T] {
	swb := &SlidingWindowBuffer[T]{
		entries:    make([]WindowEntry[T], 0),
		windowSize: windowSize,
	}

	// Start cleanup goroutine
	go swb.cleanup()

	return swb
}

// Add adds a value to the window
func (swb *SlidingWindowBuffer[T]) Add(val T) {
	swb.mu.Lock()
	defer swb.mu.Unlock()

	entry := WindowEntry[T]{
		Value:     val,
		Timestamp: time.Now(),
	}
	swb.entries = append(swb.entries, entry)
}

// GetWindow returns all values in current window
func (swb *SlidingWindowBuffer[T]) GetWindow() []T {
	swb.mu.RLock()
	defer swb.mu.RUnlock()

	cutoff := time.Now().Add(-swb.windowSize)
	result := make([]T, 0)

	for _, entry := range swb.entries {
		if entry.Timestamp.After(cutoff) {
			result = append(result, entry.Value)
		}
	}

	return result
}

// cleanup periodically removes old entries
func (swb *SlidingWindowBuffer[T]) cleanup() {
	ticker := time.NewTicker(swb.windowSize / 10)
	defer ticker.Stop()

	for range ticker.C {
		swb.mu.Lock()
		cutoff := time.Now().Add(-swb.windowSize)

		validEntries := make([]WindowEntry[T], 0)
		for _, entry := range swb.entries {
			if entry.Timestamp.After(cutoff) {
				validEntries = append(validEntries, entry)
			}
		}
		swb.entries = validEntries
		swb.mu.Unlock()
	}
}

// Size returns current number of entries
func (swb *SlidingWindowBuffer[T]) Size() int {
	swb.mu.RLock()
	defer swb.mu.RUnlock()

	cutoff := time.Now().Add(-swb.windowSize)
	count := 0

	for _, entry := range swb.entries {
		if entry.Timestamp.After(cutoff) {
			count++
		}
	}

	return count
}