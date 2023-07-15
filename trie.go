package main

type Node struct {
	end      bool
	children map[byte]*Node
}

type Trie struct {
	root *Node
}

func NewTrie() Trie {
	children := make(map[byte]*Node, 26)
	chars := &Node{false, children}
	return Trie{chars}
}

func (t *Trie) Insert(word string) {
	curr := t.root

	for i := 0; i < len(word); i++ {
		b := word[i]
		if _, exists := curr.children[b]; !exists {
			children := make(map[byte]*Node, 26)
			curr.children[b] = &Node{false, children}
		} else {
		}
		curr = curr.children[b]
	}
	curr.end = true
}

func (t *Trie) Search(word string) bool {
	curr := t.root
	for i := 0; i < len(word); i++ {
		char := word[i]
		if next, exists := curr.children[char]; !exists {
			return false
		} else {
			curr = next
		}
	}
	return curr.end
}

func (t *Trie) StartsWith(prefix string) bool {
	curr := t.root
	for i := 0; i < len(prefix); i++ {
		char := prefix[i]
		if next, exists := curr.children[char]; !exists {
			return false
		} else {
			curr = next
		}
	}
	return true
}
