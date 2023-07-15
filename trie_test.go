package main

import (
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()

	trie.Insert("apple")
	if !trie.Search("apple") {
		t.Fatal("searched for apple but could not find")
	}
	if trie.Search("app") {
		t.Fatal("did not insert the word 'app' but it was found")
	}
	trie.Insert("app")
	if !trie.Search("app") {
		t.Fatal("inserted and searched for word 'app' but could not find")
	}
}
