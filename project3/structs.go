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

// type Edge struct {
// 	pair [2]Node
// }

// func (p *Edge) sort() {
// 	if p.pair[0].ID > p.pair[1].ID {
// 		p.pair[0], p.pair[1] = p.pair[1], p.pair[0]
// 	}
// }

// func (p Edge) distance() float64 {
// 	return distance(p.pair[0], p.pair[1])
// }

type Edge struct {
	from Node
	to   Node
	cost float64
}
