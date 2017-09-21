package gorasick

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"bufio"
)

type transitionKey struct{
	state int
	symbol rune
}

type Match struct{
	location int
	text string 
}

type automaton struct{
	transitions map[int]map[rune]int
	failures map[int]int
	outputs map[int][]string
	statesCount int
}

func EmptyTrie() automaton {
	return automaton{
		transitions: make(map[int]map[rune]int),
		failures: make(map[int]int),
		statesCount: 1,
		outputs: make(map[int][]string),
	}
}

func (a *automaton) AddPrefix(prefix string) {
	var state int = 0
	var transitionsForState map[rune]int
	var ok bool
	for _, r := range prefix {
		transitionsForState, ok = a.transitions[state]
		if !ok {
			transitionsForState = make(map[rune]int)
			a.transitions[state] = transitionsForState
			state = a.statesCount
			transitionsForState[r] = state
			a.statesCount ++
			continue
		} 
		new_state, ok := transitionsForState[r]
		if ok {
			state = new_state
		} else {
			state = a.statesCount
			transitionsForState[r] = state
			a.statesCount ++
		}
	}
	a.outputs[state] = []string{prefix}
}

// goto function. Negative result is fail
func (a *automaton) g(state int, symbol rune) int {
	transitionsForState, ok := a.transitions[state]
	if !ok {
		if state == 0 { // Never fail for starting state
			return 0
		}
		return -1
	}
	res, ok := transitionsForState[symbol]
	if ok {
		return res
	} else {
		if state == 0 { // Never fail for starting state
			return 0
		}
		return -1
	}
}

func (a *automaton) buildFailures() {
	queue := make([]int, 0)
	for _, v := range a.transitions[0] {
		queue = append(queue, v)
		a.failures[v] = 0 
	}
	// fmt.Println("Starting with following failures and queue:")
	// fmt.Println(a.failures)
	// fmt.Println(queue)
	for len(queue) != 0 {
		r := queue[0]
		queue = queue[1:]

		// fmt.Printf("\nTaking %d from queue as r\n", r)

		for k, v := range a.transitions[r] {
			// fmt.Printf("%c -> %d\n", k, v)
			queue = append(queue, v)
			state := a.failures[r]
			// fmt.Printf("state is %d\n", state)
			for a.g(state, k) == -1 {
				state = a.failures[state]
				fmt.Printf("state is %d\n", state)
			}
			// fmt.Printf("No failures for %c in state %d\n", k, state)
			a.failures[v] = a.g(state, k)
			// fmt.Printf("Adding new failure for %d to %d\n", v, a.g(state, k))
			a.outputs[v] = append(a.outputs[v], a.outputs[a.failures[v]]...)
		}
	}
}

func (a *automaton) FindAll(text string) []Match {
	res := make([]Match, 0)
	state := 0
	for i, symbol := range text {
		for a.g(state, symbol) == -1 {
			state = a.failures[state]
		}
		state = a.g(state, symbol)
		for _, output := range a.outputs[state] {
			res = append(res, Match{
				location: i - len(output) + 1,
				text: output,
			})
		}
	}
	return res
}

func (a automaton) String() string {
	res := "digraph trie {\nrankdir=\"LR\";\n"
	for s1, trans := range a.transitions {
		for sym, s2 := range trans {
			res += fmt.Sprintf("\t%d -> %d [label=%c];\n", s1, s2, sym)
		}
	}
	for k, v := range a.failures {
		res += fmt.Sprintf("\t%d -> %d [label=fail; constraint=false];\n", k, v)
	}
	for k, v := range a.outputs {
		if len(v) > 0 {
			res += fmt.Sprintf(
				"\t%d -> %s [style=dotted];\n",
				k, strings.Join(v, "_"),
			)
		}
	}
	return res + "}"
}

func (a *automaton) LoadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a.AddPrefix(scanner.Text())
	}
	return scanner.Err()
}

func (a automaton) ToDotFile(filename string) {
	s := []byte(fmt.Sprintln(a))
	ioutil.WriteFile(filename, s, 0644)
}
