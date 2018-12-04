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

// Pair stores bath From one Node To another as well as the distance between
// them, MAY NOT BE NEEDED
type Pair struct {
	First  Node
	Second Node
	dist   float64
}

//needs testing inside slices
func (p *Pair) sort() {
	if p.Second.ID < p.First.ID {
		p.First, p.Second = p.Second, p.First
	}
}

func (p *Pair) distance() {
	p.dist = distance(p.First, p.Second)
}

func (p Pair) String() string {
	return fmt.Sprintf("[%v %v]", p.First, p.Second)
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
