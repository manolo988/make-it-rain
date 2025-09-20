package data_structures

import (
	"strings"
	"sync"
)

// TrieNode represents a node in the trie
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	value    interface{} // Optional value associated with the word
}

// Trie is a prefix tree for efficient string operations
type Trie struct {
	root *TrieNode
	size int
	mu   sync.RWMutex
}

// NewTrie creates a new empty trie
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
		size: 0,
	}
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[ch]
	}

	if !node.isEnd {
		t.size++
		node.isEnd = true
	}
}

// InsertWithValue adds a word with associated value
func (t *Trie) InsertWithValue(word string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[ch]
	}

	if !node.isEnd {
		t.size++
	}
	node.isEnd = true
	node.value = value
}

// Search checks if word exists in trie
func (t *Trie) Search(word string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.searchNode(word)
	return node != nil && node.isEnd
}

// SearchWithValue returns the value associated with word
func (t *Trie) SearchWithValue(word string) (interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.searchNode(word)
	if node != nil && node.isEnd {
		return node.value, true
	}
	return nil, false
}

// StartsWith checks if any word starts with prefix
func (t *Trie) StartsWith(prefix string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.searchNode(prefix) != nil
}

// searchNode helper to find node for given string
func (t *Trie) searchNode(s string) *TrieNode {
	node := t.root
	for _, ch := range s {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

// GetAllWithPrefix returns all words with given prefix
func (t *Trie) GetAllWithPrefix(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	results := []string{}
	node := t.searchNode(prefix)

	if node == nil {
		return results
	}

	// DFS to collect all words
	var dfs func(*TrieNode, string)
	dfs = func(n *TrieNode, current string) {
		if n.isEnd {
			results = append(results, current)
		}
		for ch, child := range n.children {
			dfs(child, current+string(ch))
		}
	}

	dfs(node, prefix)
	return results
}

// AutoComplete returns top N suggestions for prefix
func (t *Trie) AutoComplete(prefix string, maxSuggestions int) []string {
	allWords := t.GetAllWithPrefix(prefix)

	if len(allWords) <= maxSuggestions {
		return allWords
	}

	return allWords[:maxSuggestions]
}

// Delete removes a word from trie
func (t *Trie) Delete(word string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	var deleteHelper func(*TrieNode, string, int) bool

	deleteHelper = func(node *TrieNode, word string, index int) bool {
		if index == len(word) {
			// We've reached the end of the word
			if !node.isEnd {
				return false // Word doesn't exist
			}
			node.isEnd = false
			node.value = nil
			t.size--

			// Return true if node has no children (can be deleted)
			return len(node.children) == 0
		}

		ch := rune(word[index])
		childNode, exists := node.children[ch]
		if !exists {
			return false // Word doesn't exist
		}

		shouldDeleteChild := deleteHelper(childNode, word, index+1)

		if shouldDeleteChild {
			delete(node.children, ch)
			// Return true if current node can be deleted
			return !node.isEnd && len(node.children) == 0
		}

		return false
	}

	deleteHelper(t.root, word, 0)
	return true
}

// Size returns number of words in trie
func (t *Trie) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.size
}

// IsEmpty returns true if trie is empty
func (t *Trie) IsEmpty() bool {
	return t.Size() == 0
}

// Clear removes all words
func (t *Trie) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.root = &TrieNode{
		children: make(map[rune]*TrieNode),
	}
	t.size = 0
}

// LongestCommonPrefix finds the longest common prefix of all words
func (t *Trie) LongestCommonPrefix() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.size == 0 {
		return ""
	}

	prefix := ""
	node := t.root

	for len(node.children) == 1 && !node.isEnd {
		for ch, child := range node.children {
			prefix += string(ch)
			node = child
			break
		}
	}

	return prefix
}

// GetAll returns all words in the trie
func (t *Trie) GetAll() []string {
	return t.GetAllWithPrefix("")
}

// WildcardSearch searches with '.' as wildcard for any character
func (t *Trie) WildcardSearch(pattern string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	results := []string{}

	var search func(*TrieNode, string, int, string)
	search = func(node *TrieNode, pattern string, index int, current string) {
		if index == len(pattern) {
			if node.isEnd {
				results = append(results, current)
			}
			return
		}

		ch := pattern[index]
		if ch == '.' {
			// Wildcard - try all children
			for c, child := range node.children {
				search(child, pattern, index+1, current+string(c))
			}
		} else {
			// Regular character
			if child, exists := node.children[rune(ch)]; exists {
				search(child, pattern, index+1, current+string(ch))
			}
		}
	}

	search(t.root, pattern, 0, "")
	return results
}

// URLRouter is a specialized trie for URL routing
type URLRouter struct {
	trie *Trie
}

// NewURLRouter creates a new URL router
func NewURLRouter() *URLRouter {
	return &URLRouter{
		trie: NewTrie(),
	}
}

// AddRoute adds a URL pattern with handler
func (r *URLRouter) AddRoute(pattern string, handler interface{}) {
	// Normalize path
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	r.trie.InsertWithValue(pattern, handler)
}

// FindRoute matches a URL and returns handler
func (r *URLRouter) FindRoute(url string) (interface{}, bool) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	// Exact match
	if handler, found := r.trie.SearchWithValue(url); found {
		return handler, true
	}

	// Try prefix matching for parametric routes
	parts := strings.Split(url, "/")
	for i := len(parts); i > 0; i-- {
		prefix := strings.Join(parts[:i], "/")
		if handler, found := r.trie.SearchWithValue(prefix); found {
			return handler, true
		}
	}

	return nil, false
}

// CompressedTrie is a space-optimized trie (Radix Tree/Patricia Trie)
type CompressedTrie struct {
	root *CompressedTrieNode
	size int
	mu   sync.RWMutex
}

type CompressedTrieNode struct {
	label    string
	children map[byte]*CompressedTrieNode
	isEnd    bool
	value    interface{}
}

// NewCompressedTrie creates a new compressed trie
func NewCompressedTrie() *CompressedTrie {
	return &CompressedTrie{
		root: &CompressedTrieNode{
			children: make(map[byte]*CompressedTrieNode),
		},
		size: 0,
	}
}

// Insert adds a word to compressed trie
func (ct *CompressedTrie) Insert(word string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if len(word) == 0 {
		return
	}

	node := ct.root
	remaining := word

	for len(remaining) > 0 {
		firstChar := remaining[0]

		if child, exists := node.children[firstChar]; exists {
			// Find common prefix
			commonLen := ct.commonPrefixLength(child.label, remaining)

			if commonLen == len(child.label) {
				// child.label is a prefix of remaining
				remaining = remaining[commonLen:]
				node = child
			} else {
				// Need to split the node
				// Create new node with common part
				commonPart := child.label[:commonLen]
				childRemainder := child.label[commonLen:]
				wordRemainder := remaining[commonLen:]

				// Create new intermediate node
				newNode := &CompressedTrieNode{
					label:    commonPart,
					children: make(map[byte]*CompressedTrieNode),
					isEnd:    false,
				}

				// Update child with remainder
				child.label = childRemainder

				// Connect nodes
				node.children[firstChar] = newNode
				newNode.children[childRemainder[0]] = child

				if len(wordRemainder) == 0 {
					newNode.isEnd = true
					ct.size++
					return
				}

				// Add new branch for word remainder
				newChild := &CompressedTrieNode{
					label:    wordRemainder,
					children: make(map[byte]*CompressedTrieNode),
					isEnd:    true,
				}
				newNode.children[wordRemainder[0]] = newChild
				ct.size++
				return
			}
		} else {
			// No matching child, create new one
			newNode := &CompressedTrieNode{
				label:    remaining,
				children: make(map[byte]*CompressedTrieNode),
				isEnd:    true,
			}
			node.children[firstChar] = newNode
			ct.size++
			return
		}
	}

	if !node.isEnd {
		node.isEnd = true
		ct.size++
	}
}

// commonPrefixLength finds length of common prefix
func (ct *CompressedTrie) commonPrefixLength(s1, s2 string) int {
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return minLen
}

// Search checks if word exists
func (ct *CompressedTrie) Search(word string) bool {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	node := ct.root
	remaining := word

	for len(remaining) > 0 {
		firstChar := remaining[0]

		if child, exists := node.children[firstChar]; exists {
			if strings.HasPrefix(remaining, child.label) {
				remaining = remaining[len(child.label):]
				node = child
			} else {
				return false
			}
		} else {
			return false
		}
	}

	return node.isEnd
}