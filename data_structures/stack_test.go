package data_structures

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack[int]()

	// Test empty stack
	if !s.IsEmpty() {
		t.Error("New stack should be empty")
	}

	if s.Size() != 0 {
		t.Errorf("Expected size 0, got %d", s.Size())
	}

	// Test push
	s.Push(10)
	s.Push(20)
	s.Push(30)

	if s.Size() != 3 {
		t.Errorf("Expected size 3, got %d", s.Size())
	}

	// Test peek
	if val, ok := s.Peek(); !ok || val != 30 {
		t.Errorf("Expected peek 30, got %d", val)
	}

	// Test pop
	if val, ok := s.Pop(); !ok || val != 30 {
		t.Errorf("Expected pop 30, got %d", val)
	}

	if val, ok := s.Pop(); !ok || val != 20 {
		t.Errorf("Expected pop 20, got %d", val)
	}

	// Test ToSlice
	s.Push(40)
	slice := s.ToSlice()
	expected := []int{10, 40}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v, got %v", expected, slice)
	}

	// Test clear
	s.Clear()
	if !s.IsEmpty() {
		t.Error("Stack should be empty after clear")
	}

	// Test pop on empty stack
	if _, ok := s.Pop(); ok {
		t.Error("Pop on empty stack should return false")
	}
}

func TestMinStack(t *testing.T) {
	ms := NewMinStack[int]()

	ms.Push(3)
	ms.Push(5)
	if min, _ := ms.GetMin(); min != 3 {
		t.Errorf("Expected min 3, got %d", min)
	}

	ms.Push(2)
	ms.Push(1)
	if min, _ := ms.GetMin(); min != 1 {
		t.Errorf("Expected min 1, got %d", min)
	}

	ms.Pop()
	if min, _ := ms.GetMin(); min != 2 {
		t.Errorf("Expected min 2 after pop, got %d", min)
	}

	ms.Pop()
	if min, _ := ms.GetMin(); min != 3 {
		t.Errorf("Expected min 3 after pop, got %d", min)
	}

	// Test duplicate minimums
	ms.Push(1)
	ms.Push(1)
	ms.Pop()
	if min, _ := ms.GetMin(); min != 1 {
		t.Errorf("Expected min 1 with duplicates, got %d", min)
	}
}

func TestMaxStack(t *testing.T) {
	ms := NewMaxStack[int]()

	ms.Push(3)
	ms.Push(5)
	if max, _ := ms.GetMax(); max != 5 {
		t.Errorf("Expected max 5, got %d", max)
	}

	ms.Push(7)
	ms.Push(2)
	if max, _ := ms.GetMax(); max != 7 {
		t.Errorf("Expected max 7, got %d", max)
	}

	ms.Pop() // Remove 2
	ms.Pop() // Remove 7
	if max, _ := ms.GetMax(); max != 5 {
		t.Errorf("Expected max 5 after pops, got %d", max)
	}
}

func TestMonotonicStack(t *testing.T) {
	// Test increasing monotonic stack
	incStack := NewMonotonicStack[int](true)

	// Push 3 - stack: [3]
	popped := incStack.Push(3)
	if len(popped) != 0 {
		t.Error("Should not pop any elements")
	}

	// Push 1 - stack: [1] (3 is popped)
	popped = incStack.Push(1)
	if len(popped) != 1 || popped[0] != 3 {
		t.Errorf("Should pop 3, got %v", popped)
	}

	// Push 4 - stack: [1, 4]
	popped = incStack.Push(4)
	if len(popped) != 0 {
		t.Error("Should not pop any elements")
	}

	// Push 2 - stack: [1, 2] (4 is popped)
	popped = incStack.Push(2)
	if len(popped) != 1 || popped[0] != 4 {
		t.Errorf("Should pop 4, got %v", popped)
	}

	expected := []int{1, 2}
	if !reflect.DeepEqual(incStack.ToSlice(), expected) {
		t.Errorf("Expected stack %v, got %v", expected, incStack.ToSlice())
	}

	// Test decreasing monotonic stack
	decStack := NewMonotonicStack[int](false)

	// Push 1, 4, 3, 2
	decStack.Push(1)
	decStack.Push(4)
	popped = decStack.Push(3)
	if len(popped) != 0 {
		t.Error("Should not pop any elements")
	}

	popped = decStack.Push(2)
	if len(popped) != 0 {
		t.Error("Should not pop any elements")
	}

	// Stack should be [4, 3, 2] in decreasing order
	expected = []int{4, 3, 2}
	if !reflect.DeepEqual(decStack.ToSlice(), expected) {
		t.Errorf("Expected stack %v, got %v", expected, decStack.ToSlice())
	}

	// Push 5 - should pop all elements smaller than 5
	popped = decStack.Push(5)
	if len(popped) != 3 {
		t.Errorf("Should pop 3 elements, got %d", len(popped))
	}

	if decStack.Size() != 1 {
		t.Errorf("Stack should have 1 element, got %d", decStack.Size())
	}
}

func TestStackWithStrings(t *testing.T) {
	s := NewStack[string]()

	s.Push("hello")
	s.Push("world")

	if val, _ := s.Peek(); val != "world" {
		t.Errorf("Expected 'world', got %s", val)
	}

	if val, _ := s.Pop(); val != "world" {
		t.Errorf("Expected 'world', got %s", val)
	}

	if val, _ := s.Pop(); val != "hello" {
		t.Errorf("Expected 'hello', got %s", val)
	}

	if !s.IsEmpty() {
		t.Error("Stack should be empty")
	}
}

func TestStackCapacity(t *testing.T) {
	s := NewStack[int]()

	// Push many elements to test dynamic growth
	n := 1000
	for i := 0; i < n; i++ {
		s.Push(i)
	}

	if s.Size() != n {
		t.Errorf("Expected size %d, got %d", n, s.Size())
	}

	// Pop all elements
	for i := n - 1; i >= 0; i-- {
		val, ok := s.Pop()
		if !ok || val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all elements")
	}
}