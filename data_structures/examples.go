package data_structures

// Common Interview Problems Using Heaps

// 1. Kth Largest Element in Array
func FindKthLargest(nums []int, k int) int {
	minHeap := NewMinHeap[int]()

	for _, num := range nums {
		minHeap.Push(num)
		if minHeap.Size() > k {
			minHeap.Pop()
		}
	}

	result, _ := minHeap.Peek()
	return result
}

// 2. Merge K Sorted Lists
type ListNode struct {
	Val  int
	Next *ListNode
}

func MergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	type NodeWithIndex struct {
		node  *ListNode
		index int
	}

	h := NewHeapWithComparator(func(a, b NodeWithIndex) bool {
		return a.node.Val < b.node.Val
	})

	// Add first node from each list
	for i, list := range lists {
		if list != nil {
			h.Push(NodeWithIndex{node: list, index: i})
		}
	}

	dummy := &ListNode{}
	current := dummy

	for !h.IsEmpty() {
		item, _ := h.Pop()
		current.Next = item.node
		current = current.Next

		if item.node.Next != nil {
			h.Push(NodeWithIndex{
				node:  item.node.Next,
				index: item.index,
			})
		}
	}

	return dummy.Next
}

// 3. Top K Frequent Elements
func TopKFrequent(nums []int, k int) []int {
	// Count frequencies
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	type FreqPair struct {
		num   int
		count int
	}

	// Min heap of size k
	h := NewHeapWithComparator(func(a, b FreqPair) bool {
		return a.count < b.count
	})

	for num, count := range freq {
		h.Push(FreqPair{num: num, count: count})
		if h.Size() > k {
			h.Pop()
		}
	}

	result := make([]int, 0, k)
	for !h.IsEmpty() {
		item, _ := h.Pop()
		result = append(result, item.num)
	}

	return result
}

// 4. Task Scheduler (Leetcode 621)
func LeastInterval(tasks []byte, n int) int {
	if n == 0 {
		return len(tasks)
	}

	// Count task frequencies
	freq := make(map[byte]int)
	for _, task := range tasks {
		freq[task]++
	}

	// Max heap based on frequency
	maxHeap := NewMaxHeap[int]()
	for _, count := range freq {
		maxHeap.Push(count)
	}

	intervals := 0

	for !maxHeap.IsEmpty() {
		temp := make([]int, 0)

		// Try to execute n+1 different tasks
		for i := 0; i <= n; i++ {
			if !maxHeap.IsEmpty() {
				count, _ := maxHeap.Pop()
				if count > 1 {
					temp = append(temp, count-1)
				}
			}
			intervals++

			if maxHeap.IsEmpty() && len(temp) == 0 {
				break
			}
		}

		// Add back tasks that still need execution
		for _, count := range temp {
			maxHeap.Push(count)
		}
	}

	return intervals
}

// 5. Meeting Rooms II (minimum meeting rooms needed)
type Interval struct {
	Start int
	End   int
}

func MinMeetingRooms(intervals []Interval) int {
	if len(intervals) == 0 {
		return 0
	}

	// Sort by start time
	sortedIntervals := make([]Interval, len(intervals))
	copy(sortedIntervals, intervals)
	// Note: In real implementation, sort intervals by start time

	// Min heap to track end times
	minHeap := NewMinHeap[int]()
	minHeap.Push(sortedIntervals[0].End)

	for i := 1; i < len(sortedIntervals); i++ {
		earliest, _ := minHeap.Peek()

		if sortedIntervals[i].Start >= earliest {
			// Reuse the room
			minHeap.Pop()
		}

		minHeap.Push(sortedIntervals[i].End)
	}

	return minHeap.Size()
}

// 6. Dijkstra's Shortest Path Algorithm
type Edge struct {
	To     int
	Weight int
}

func Dijkstra(graph map[int][]Edge, start int, end int) int {
	type State struct {
		node int
		dist int
	}

	minHeap := NewHeapWithComparator(func(a, b State) bool {
		return a.dist < b.dist
	})

	distances := make(map[int]int)
	visited := make(map[int]bool)

	minHeap.Push(State{node: start, dist: 0})
	distances[start] = 0

	for !minHeap.IsEmpty() {
		current, _ := minHeap.Pop()

		if current.node == end {
			return current.dist
		}

		if visited[current.node] {
			continue
		}
		visited[current.node] = true

		for _, edge := range graph[current.node] {
			newDist := current.dist + edge.Weight

			if oldDist, exists := distances[edge.To]; !exists || newDist < oldDist {
				distances[edge.To] = newDist
				minHeap.Push(State{node: edge.To, dist: newDist})
			}
		}
	}

	return -1 // Path not found
}

// 7. Sliding Window Maximum
func MaxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}

	type IndexedValue struct {
		value int
		index int
	}

	maxHeap := NewHeapWithComparator(func(a, b IndexedValue) bool {
		return a.value > b.value
	})

	result := make([]int, 0, len(nums)-k+1)

	// Initialize window
	for i := 0; i < k && i < len(nums); i++ {
		maxHeap.Push(IndexedValue{value: nums[i], index: i})
	}

	if maxHeap.Size() > 0 {
		val, _ := maxHeap.Peek()
		result = append(result, val.value)
	}

	// Slide window
	for i := k; i < len(nums); i++ {
		maxHeap.Push(IndexedValue{value: nums[i], index: i})

		// Remove elements outside window
		for !maxHeap.IsEmpty() {
			top, _ := maxHeap.Peek()
			if top.index <= i-k {
				maxHeap.Pop()
			} else {
				break
			}
		}

		if !maxHeap.IsEmpty() {
			val, _ := maxHeap.Peek()
			result = append(result, val.value)
		}
	}

	return result
}

// 8. Kth Smallest Element in Sorted Matrix
func KthSmallestInMatrix(matrix [][]int, k int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}

	type Cell struct {
		val int
		row int
		col int
	}

	n := len(matrix)
	minHeap := NewHeapWithComparator(func(a, b Cell) bool {
		return a.val < b.val
	})

	// Add first element of each row
	for i := 0; i < n && i < k; i++ {
		minHeap.Push(Cell{val: matrix[i][0], row: i, col: 0})
	}

	var result Cell
	for i := 0; i < k && !minHeap.IsEmpty(); i++ {
		result, _ = minHeap.Pop()

		// Add next element from same row
		if result.col+1 < n {
			minHeap.Push(Cell{
				val: matrix[result.row][result.col+1],
				row: result.row,
				col: result.col + 1,
			})
		}
	}

	return result.val
}