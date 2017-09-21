package gorasick

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

func SheHisHers() automaton {
	trie := EmptyTrie()
	trie.AddPrefix("he")
	trie.AddPrefix("she")
	trie.AddPrefix("his")
	trie.AddPrefix("hers")
	trie.ToDotFile("without_failures.dot")
	trie.buildFailures()
	trie.ToDotFile("hello.dot")
	// dot -Tsvg hello.dot | display # to view that file
	return trie
}

func ExampleSearch() {
	trie := SheHisHers()
	text := "shers his hello"
	fmt.Println(text)
	for _, res := range trie.FindAll(text) {
		for i := 0; i < res.location; i++ {
			fmt.Printf(" ")
		}
		fmt.Printf("%s\n", res.text)
	}
	// Output:
	// shers his hello
	// she
	//  he
	//  hers
	//       his
	//           he
}

func TestFailures(t *testing.T) {
	t.Run("hers he she his", func(t *testing.T) {
		trie := SheHisHers()

		var f = []struct {s1 int; s2 int}{
			{1, 0},
			{2, 0},
			{3, 0},
			{4, 1},
			{5, 2},
			{6, 0},
			{7, 3},
			{8, 0},
			{9, 3},
		}

		for _, tc := range f {
			assert.Equal(t, trie.failures[tc.s1], tc.s2)
		}
	})
}
