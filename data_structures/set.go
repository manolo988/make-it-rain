package data_structures

import (
	"sync"
)

// Set is a generic hash set implementation
type Set[T comparable] struct {
	data map[T]struct{}
	mu   sync.RWMutex
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
	}
}

// NewSetFrom creates a set from slice
func NewSetFrom[T comparable](items []T) *Set[T] {
	s := NewSet[T]()
	for _, item := range items {
		s.Add(item)
	}
	return s
}

// Add adds an element to the set
func (s *Set[T]) Add(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[item]; exists {
		return false
	}
	s.data[item] = struct{}{}
	return true
}

// AddAll adds multiple elements
func (s *Set[T]) AddAll(items ...T) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	added := 0
	for _, item := range items {
		if _, exists := s.data[item]; !exists {
			s.data[item] = struct{}{}
			added++
		}
	}
	return added
}

// Remove removes an element from the set
func (s *Set[T]) Remove(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[item]; exists {
		delete(s.data, item)
		return true
	}
	return false
}

// Contains checks if element exists
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[item]
	return exists
}

// ContainsAll checks if all elements exist
func (s *Set[T]) ContainsAll(items ...T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range items {
		if _, exists := s.data[item]; !exists {
			return false
		}
	}
	return true
}

// ContainsAny checks if any element exists
func (s *Set[T]) ContainsAny(items ...T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range items {
		if _, exists := s.data[item]; exists {
			return true
		}
	}
	return false
}

// Size returns the number of elements
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// IsEmpty returns true if set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear removes all elements
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[T]struct{})
}

// ToSlice returns all elements as slice
func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]T, 0, len(s.data))
	for item := range s.data {
		result = append(result, item)
	}
	return result
}

// Clone creates a copy of the set
func (s *Set[T]) Clone() *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newSet := NewSet[T]()
	for item := range s.data {
		newSet.data[item] = struct{}{}
	}
	return newSet
}

// Union returns a new set with elements from both sets
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSet[T]()
	for item := range s.data {
		result.data[item] = struct{}{}
	}
	for item := range other.data {
		result.data[item] = struct{}{}
	}
	return result
}

// Intersection returns a new set with common elements
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSet[T]()

	// Iterate over smaller set for efficiency
	var smaller, larger map[T]struct{}
	if len(s.data) <= len(other.data) {
		smaller = s.data
		larger = other.data
	} else {
		smaller = other.data
		larger = s.data
	}

	for item := range smaller {
		if _, exists := larger[item]; exists {
			result.data[item] = struct{}{}
		}
	}
	return result
}

// Difference returns elements in s but not in other
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSet[T]()
	for item := range s.data {
		if _, exists := other.data[item]; !exists {
			result.data[item] = struct{}{}
		}
	}
	return result
}

// SymmetricDifference returns elements in either set but not both
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSet[T]()

	// Add elements only in s
	for item := range s.data {
		if _, exists := other.data[item]; !exists {
			result.data[item] = struct{}{}
		}
	}

	// Add elements only in other
	for item := range other.data {
		if _, exists := s.data[item]; !exists {
			result.data[item] = struct{}{}
		}
	}

	return result
}

// IsSubset checks if s is a subset of other
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	if len(s.data) > len(other.data) {
		return false
	}

	for item := range s.data {
		if _, exists := other.data[item]; !exists {
			return false
		}
	}
	return true
}

// IsSuperset checks if s is a superset of other
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

// IsDisjoint checks if sets have no common elements
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	// Check smaller set against larger
	var smaller, larger map[T]struct{}
	if len(s.data) <= len(other.data) {
		smaller = s.data
		larger = other.data
	} else {
		smaller = other.data
		larger = s.data
	}

	for item := range smaller {
		if _, exists := larger[item]; exists {
			return false
		}
	}
	return true
}

// Equals checks if sets contain the same elements
func (s *Set[T]) Equals(other *Set[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	if len(s.data) != len(other.data) {
		return false
	}

	for item := range s.data {
		if _, exists := other.data[item]; !exists {
			return false
		}
	}
	return true
}

// ForEach applies function to each element
func (s *Set[T]) ForEach(fn func(T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for item := range s.data {
		fn(item)
	}
}

// Filter returns a new set with elements that pass the test
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := NewSet[T]()
	for item := range s.data {
		if predicate(item) {
			result.data[item] = struct{}{}
		}
	}
	return result
}

// Pop removes and returns an arbitrary element
func (s *Set[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var zero T
	for item := range s.data {
		delete(s.data, item)
		return item, true
	}
	return zero, false
}