package data_structures

import (
	"strconv"
	"strings"
)

// Common Stack Interview Problems

// 1. Valid Parentheses
func IsValid(s string) bool {
	stack := NewStack[rune]()
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, ch := range s {
		if ch == '(' || ch == '{' || ch == '[' {
			stack.Push(ch)
		} else if open, exists := pairs[ch]; exists {
			if top, ok := stack.Pop(); !ok || top != open {
				return false
			}
		}
	}

	return stack.IsEmpty()
}

// 2. Evaluate Reverse Polish Notation
func EvalRPN(tokens []string) int {
	stack := NewStack[int]()

	for _, token := range tokens {
		if token == "+" || token == "-" || token == "*" || token == "/" {
			b, _ := stack.Pop()
			a, _ := stack.Pop()

			var result int
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				result = a / b
			}
			stack.Push(result)
		} else {
			num, _ := strconv.Atoi(token)
			stack.Push(num)
		}
	}

	result, _ := stack.Pop()
	return result
}

// 3. Daily Temperatures - Days to wait for warmer temperature
func DailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	result := make([]int, n)
	stack := NewStack[int]() // Stack of indices

	for i := 0; i < n; i++ {
		for !stack.IsEmpty() {
			topIdx, _ := stack.Peek()
			if temperatures[i] > temperatures[topIdx] {
				stack.Pop()
				result[topIdx] = i - topIdx
			} else {
				break
			}
		}
		stack.Push(i)
	}

	return result
}

// 4. Next Greater Element
func NextGreaterElement(nums []int) []int {
	n := len(nums)
	result := make([]int, n)
	for i := range result {
		result[i] = -1
	}

	stack := NewStack[int]() // Stack of indices

	for i := 0; i < n; i++ {
		for !stack.IsEmpty() {
			topIdx, _ := stack.Peek()
			if nums[i] > nums[topIdx] {
				stack.Pop()
				result[topIdx] = nums[i]
			} else {
				break
			}
		}
		stack.Push(i)
	}

	return result
}

// 5. Largest Rectangle in Histogram
func LargestRectangleArea(heights []int) int {
	stack := NewStack[int]()
	maxArea := 0
	n := len(heights)

	for i := 0; i < n; i++ {
		for !stack.IsEmpty() {
			topIdx, _ := stack.Peek()
			if heights[i] < heights[topIdx] {
				stack.Pop()
				height := heights[topIdx]

				width := i
				if !stack.IsEmpty() {
					prevIdx, _ := stack.Peek()
					width = i - prevIdx - 1
				}

				area := height * width
				if area > maxArea {
					maxArea = area
				}
			} else {
				break
			}
		}
		stack.Push(i)
	}

	// Process remaining bars
	for !stack.IsEmpty() {
		topIdx, _ := stack.Pop()
		height := heights[topIdx]

		width := n
		if !stack.IsEmpty() {
			prevIdx, _ := stack.Peek()
			width = n - prevIdx - 1
		}

		area := height * width
		if area > maxArea {
			maxArea = area
		}
	}

	return maxArea
}

// 6. Basic Calculator (supports +, -, (, ))
func Calculate(s string) int {
	stack := NewStack[int]()
	result := 0
	number := 0
	sign := 1

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch >= '0' && ch <= '9' {
			number = number*10 + int(ch-'0')
		} else if ch == '+' {
			result += sign * number
			number = 0
			sign = 1
		} else if ch == '-' {
			result += sign * number
			number = 0
			sign = -1
		} else if ch == '(' {
			stack.Push(result)
			stack.Push(sign)
			result = 0
			sign = 1
		} else if ch == ')' {
			result += sign * number
			number = 0

			prevSign, _ := stack.Pop()
			prevResult, _ := stack.Pop()
			result = prevResult + prevSign*result
		}
	}

	result += sign * number
	return result
}

// 7. Remove K Digits to get smallest number
func RemoveKDigits(num string, k int) string {
	if k >= len(num) {
		return "0"
	}

	stack := NewStack[byte]()

	for i := 0; i < len(num); i++ {
		digit := num[i]

		// Remove larger digits from stack
		for k > 0 && !stack.IsEmpty() {
			top, _ := stack.Peek()
			if top > digit {
				stack.Pop()
				k--
			} else {
				break
			}
		}

		stack.Push(digit)
	}

	// Remove remaining k digits from end
	for k > 0 {
		stack.Pop()
		k--
	}

	// Build result
	result := stack.ToSlice()

	// Remove leading zeros
	start := 0
	for start < len(result) && result[start] == '0' {
		start++
	}

	if start == len(result) {
		return "0"
	}

	return string(result[start:])
}

// 8. Trapping Rain Water
func Trap(height []int) int {
	stack := NewStack[int]() // Stack of indices
	water := 0

	for i := 0; i < len(height); i++ {
		for !stack.IsEmpty() {
			topIdx, _ := stack.Peek()
			if height[i] > height[topIdx] {
				stack.Pop()

				if stack.IsEmpty() {
					break
				}

				prevIdx, _ := stack.Peek()
				minHeight := min(height[prevIdx], height[i]) - height[topIdx]
				width := i - prevIdx - 1
				water += minHeight * width
			} else {
				break
			}
		}
		stack.Push(i)
	}

	return water
}

// 9. Simplify Path (Unix-style)
func SimplifyPath(path string) string {
	stack := NewStack[string]()
	parts := strings.Split(path, "/")

	for _, part := range parts {
		if part == "" || part == "." {
			continue
		} else if part == ".." {
			stack.Pop()
		} else {
			stack.Push(part)
		}
	}

	if stack.IsEmpty() {
		return "/"
	}

	// Build result
	result := ""
	for _, dir := range stack.ToSlice() {
		result += "/" + dir
	}

	return result
}

// 10. Min Stack Problem - Get minimum in O(1)
type MinStackProblem struct {
	stack    *Stack[int]
	minStack *Stack[int]
}

func NewMinStackProblem() *MinStackProblem {
	return &MinStackProblem{
		stack:    NewStack[int](),
		minStack: NewStack[int](),
	}
}

func (ms *MinStackProblem) Push(val int) {
	ms.stack.Push(val)

	minVal := val
	if !ms.minStack.IsEmpty() {
		currentMin, _ := ms.minStack.Peek()
		if currentMin < val {
			minVal = currentMin
		}
	}
	ms.minStack.Push(minVal)
}

func (ms *MinStackProblem) Pop() {
	ms.stack.Pop()
	ms.minStack.Pop()
}

func (ms *MinStackProblem) Top() int {
	val, _ := ms.stack.Peek()
	return val
}

func (ms *MinStackProblem) GetMin() int {
	val, _ := ms.minStack.Peek()
	return val
}

// 11. Decode String (e.g., "3[a2[c]]" = "accaccacc")
func DecodeString(s string) string {
	countStack := NewStack[int]()
	stringStack := NewStack[string]()
	currentString := ""
	k := 0

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch >= '0' && ch <= '9' {
			k = k*10 + int(ch-'0')
		} else if ch == '[' {
			countStack.Push(k)
			stringStack.Push(currentString)
			currentString = ""
			k = 0
		} else if ch == ']' {
			count, _ := countStack.Pop()
			prevString, _ := stringStack.Pop()

			repeated := ""
			for j := 0; j < count; j++ {
				repeated += currentString
			}
			currentString = prevString + repeated
		} else {
			currentString += string(ch)
		}
	}

	return currentString
}

// 12. Asteroid Collision
func AsteroidCollision(asteroids []int) []int {
	stack := NewStack[int]()

	for _, asteroid := range asteroids {
		alive := true

		for !stack.IsEmpty() && asteroid < 0 {
			top, _ := stack.Peek()
			if top < 0 {
				break // Both moving left, no collision
			}

			// Collision occurs
			if abs(top) < abs(asteroid) {
				stack.Pop() // Top asteroid explodes
				continue
			} else if abs(top) == abs(asteroid) {
				stack.Pop() // Both explode
				alive = false
				break
			} else {
				alive = false // Current asteroid explodes
				break
			}
		}

		if alive {
			stack.Push(asteroid)
		}
	}

	return stack.ToSlice()
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}