package data_structures

import (
	"testing"
)

func TestMinHeap(t *testing.T) {
	h := NewMinHeap[int]()

	// Test push and size
	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(1)
	h.Push(9)

	if h.Size() != 5 {
		t.Errorf("Expected size 5, got %d", h.Size())
	}

	// Test peek
	if val, _ := h.Peek(); val != 1 {
		t.Errorf("Expected min 1, got %d", val)
	}

	// Test pop in order
	expected := []int{1, 3, 5, 7, 9}
	for i, want := range expected {
		got, ok := h.Pop()
		if !ok {
			t.Errorf("Pop failed at index %d", i)
		}
		if got != want {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}

	// Test empty heap
	if !h.IsEmpty() {
		t.Error("Heap should be empty")
	}
}

func TestMaxHeap(t *testing.T) {
	h := NewMaxHeap[int]()

	values := []int{5, 3, 7, 1, 9, 4, 6}
	for _, v := range values {
		h.Push(v)
	}

	// Should pop in descending order
	expected := []int{9, 7, 6, 5, 4, 3, 1}
	for _, want := range expected {
		got, ok := h.Pop()
		if !ok {
			t.Error("Pop failed")
		}
		if got != want {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}
}

func TestBuildHeap(t *testing.T) {
	data := []int{9, 5, 6, 2, 3, 7, 1, 4, 8}
	h := BuildHeap(data, func(a, b int) bool { return a < b })

	if h.Size() != len(data) {
		t.Errorf("Expected size %d, got %d", len(data), h.Size())
	}

	// Pop all elements, should be in ascending order
	prev := -1
	for !h.IsEmpty() {
		val, _ := h.Pop()
		if val < prev {
			t.Errorf("Heap order violated: %d < %d", val, prev)
		}
		prev = val
	}
}

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[string]()

	// Add tasks with priorities
	pq.Push("Low priority task", 10)
	pq.Push("High priority task", 1)
	pq.Push("Medium priority task", 5)

	// Should pop in priority order (lowest number = highest priority)
	if task, _ := pq.Pop(); task != "High priority task" {
		t.Errorf("Expected high priority task, got %s", task)
	}

	if task, _ := pq.Pop(); task != "Medium priority task" {
		t.Errorf("Expected medium priority task, got %s", task)
	}

	if task, _ := pq.Pop(); task != "Low priority task" {
		t.Errorf("Expected low priority task, got %s", task)
	}
}

func TestMaxPriorityQueue(t *testing.T) {
	pq := NewMaxPriorityQueue[string]()

	pq.Push("Task A", 1)
	pq.Push("Task B", 5)
	pq.Push("Task C", 3)

	// Should pop in descending priority order
	if task, _ := pq.Pop(); task != "Task B" {
		t.Errorf("Expected Task B, got %s", task)
	}
}

func TestIndexedHeap(t *testing.T) {
	h := NewIndexedMinHeap[int]()

	// Test push and contains
	h.Push(5)
	h.Push(3)
	h.Push(7)

	if !h.Contains(3) {
		t.Error("Should contain 3")
	}

	if h.Contains(10) {
		t.Error("Should not contain 10")
	}

	// Test update
	if !h.Update(5, 1) {
		t.Error("Update should succeed")
	}

	// 1 should now be the minimum
	if val, _ := h.Pop(); val != 1 {
		t.Errorf("Expected 1 after update, got %d", val)
	}

	// Test delete
	h.Push(4)
	h.Push(6)
	if !h.Delete(4) {
		t.Error("Delete should succeed")
	}

	if h.Contains(4) {
		t.Error("Should not contain 4 after delete")
	}

	// Remaining elements should still maintain heap property
	if val, _ := h.Pop(); val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}
}

func TestCustomComparator(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Create heap sorted by age
	h := NewHeapWithComparator(func(a, b Person) bool {
		return a.Age < b.Age
	})

	h.Push(Person{"Alice", 30})
	h.Push(Person{"Bob", 25})
	h.Push(Person{"Charlie", 35})

	// Should pop in age order
	if person, _ := h.Pop(); person.Name != "Bob" {
		t.Errorf("Expected Bob (youngest), got %s", person.Name)
	}
}

// Example: Find K Largest Elements (common interview problem)
func TestFindKLargest(t *testing.T) {
	findKLargest := func(nums []int, k int) []int {
		minHeap := NewMinHeap[int]()

		for _, num := range nums {
			minHeap.Push(num)
			if minHeap.Size() > k {
				minHeap.Pop()
			}
		}

		result := make([]int, 0, k)
		for !minHeap.IsEmpty() {
			val, _ := minHeap.Pop()
			result = append(result, val)
		}
		return result
	}

	nums := []int{3, 2, 1, 5, 6, 4}
	k := 3
	result := findKLargest(nums, k)

	// Should contain 4, 5, 6 (in any order)
	expected := map[int]bool{4: true, 5: true, 6: true}
	for _, v := range result {
		if !expected[v] {
			t.Errorf("Unexpected value in k largest: %d", v)
		}
	}
}

// Example: Median Finder (using two heaps)
type MedianFinder struct {
	maxHeap *Heap[int] // For smaller half
	minHeap *Heap[int] // For larger half
}

func NewMedianFinder() *MedianFinder {
	return &MedianFinder{
		maxHeap: NewMaxHeap[int](),
		minHeap: NewMinHeap[int](),
	}
}

func (mf *MedianFinder) AddNum(num int) {
	mf.maxHeap.Push(num)

	// Move max of smaller half to larger half
	if val, ok := mf.maxHeap.Pop(); ok {
		mf.minHeap.Push(val)
	}

	// Balance the heaps
	if mf.minHeap.Size() > mf.maxHeap.Size()+1 {
		if val, ok := mf.minHeap.Pop(); ok {
			mf.maxHeap.Push(val)
		}
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	if mf.minHeap.Size() > mf.maxHeap.Size() {
		val, _ := mf.minHeap.Peek()
		return float64(val)
	}
	val1, _ := mf.maxHeap.Peek()
	val2, _ := mf.minHeap.Peek()
	return float64(val1+val2) / 2.0
}

func TestMedianFinder(t *testing.T) {
	mf := NewMedianFinder()

	mf.AddNum(1)
	mf.AddNum(2)
	if median := mf.FindMedian(); median != 1.5 {
		t.Errorf("Expected median 1.5, got %f", median)
	}

	mf.AddNum(3)
	if median := mf.FindMedian(); median != 2.0 {
		t.Errorf("Expected median 2.0, got %f", median)
	}
}