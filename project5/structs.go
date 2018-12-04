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
	// return fmt.Sprintf("%v ", *t.path)
}

//Fitness contains a population along with its fitness score
// type Fitness struct {
// 	trip    *Trip
// 	fitness float64
// }
type Fitness struct {
	tripID  int
	fitness float64
}

// func (f Fitness) String() string {
// ID// }
func (f Fitness) String() string {
	return fmt.Sprintf("%v %v", f.tripID, f.fitness)
}

//ByFitness uses sort.Sort to sort by fitness score, the higher the fitness the
//lower the trip cost therefore the Less function is actually a Greater-Than
//function
type ByFitness []Fitness

func (b ByFitness) Len() int           { return len(b) }
func (b ByFitness) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByFitness) Less(i, j int) bool { return b[i].fitness > b[j].fitness }

type Generation struct {
	generation  int
	bestCost    float64
	averageCost float64
	worstCost   float64
}
