package main

import (
	"fmt"
)

//structs and their methods. some may not be needed.

// Node struct which holds the "city's" ID, X & Y coordinates.
type Node struct {
	ID int
	X  float64
	Y  float64
}

func (n Node) String() string {
	return fmt.Sprintf("%v", n.ID)
}

// Trip stores the total cost of trip as well as the total cost of the trip.
type Trip struct {
	path *[]Node
	cost float64
}

// stringer method for Trip struct
func (t Trip) String() string {
	return fmt.Sprintf("%v %v", *t.path, t.cost)
}

// DGraph is a representation of a directed graph.
type DGraph struct {
	node    Node
	visited int
	dist    int
	prev    *DGraph
	to      []*DGraph
}

// Colors for nodes
const (
	//White, Gray, and Black repreenting: Unvisited, Discovered, and Visited
	WHITE = iota
	GRAY
	BLACK
)

func (d DGraph) String() string {
	// return fmt.Sprintf("%v status: %v prev: %p\n", d.node, d.visited, d.prev)
	return fmt.Sprintf("%v", d.node)
}

//STACK/QUEUE implementation from: https://stackoverflow.com/questions/28541609/looking-for-reasonable-stack-implementation-in-golang
type stack []*DGraph

func (s stack) Push(v *DGraph) stack {
	return append(s, v)
}

func (s stack) Len() int {
	return len(s)
}

func (s stack) Pop() (stack, *DGraph) {
	// FIXME: What do we do if the stack is empty, though?

	l := len(s)
	return s[:l-1], s[l-1]
}

//not needed, may just use Push()
func (s stack) Queue(v *DGraph) stack {
	return s.Push(v)
}

func (s stack) Dequeue() (stack, *DGraph) {
	return s[1:], s[0]
}
