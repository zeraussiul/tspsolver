package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

//based slightly on: https://paginas.fe.up.pt/~mac/ensino/docs/OR/HowToSolveIt/ConstructiveHeuristicsForTheTSP.pdf
func greedy(nn []Node) *Trip {
	//cities hold cities to be processed into trip
	cities := make([]Node, len(nn))
	copy(cities, nn)
	// cities := nn
	//trip holds the edges/partial tour
	trip := []Node{}

	//Initialization: starting with random city, then find nearest city, add to
	//trip slice and remove from cities slice, this build the partial tour (i,
	//j)
	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(len(nn) - 1)

	var end int
	trip = append(trip, cities[start])
	fmt.Println("start:", cities[start])
	last := Node{cities[start].ID, cities[start].X, cities[start].Y}
	//remove starting city from the cities slice
	cities = append(cities[:start], cities[start+1:]...)

	//build the first edge
	shortest := 0.0
	for i, v := range cities {
		dist := distance(cities[start], v)
		if start != i && dist > shortest {
			shortest = dist
			end = i
		}
	}

	trip = append(trip, cities[end])
	//complete the trip by appending starting city again
	trip = append(trip, last)

	//remove second city from cities slice
	cities = append(cities[:end], cities[end+1:]...)

	//Selection: find cities k and j (j belonging to tour (trip), while k not belonging to tour (cities)) for which cost(k,j) is minimized
	for len(cities) != 0 {
		//insert variable tracks insertion point between insert and insert + 1 to be inserted into trip as an edge.
		var insert, city int
		var cost float64
		shortest = math.MaxFloat64

		for j := 0; j < len(trip)-1; j++ {
			for h, k := range cities {
				cost = distToEdge(k, trip[j], trip[j+1])
				if cost < shortest {
					shortest = cost
					insert = j
					city = h
				}
			}
		}

		//shortest found, insert into trip as an edge, between "insert" and
		//"insert+1", then delete from the cities slice. repeat until all cities
		//have been removed from the cities slice and added to the trip slice.
		trip = append(trip, Node{})
		copy(trip[insert+2:], trip[insert+1:])
		trip[insert+1] = cities[city]

		cities = append(cities[:city], cities[city+1:]...)

		//uncomment to visualize path finding.
		// drawP(trip, cities)
		// time.Sleep(1500 * time.Millisecond)

	}

	fmt.Println("trip:", trip, "dist:", pathDistance(trip))

	final := &Trip{&trip, pathDistance(trip)}
	fmt.Println("inside greedy:", nn)

	return final
}

// distToEdge calculates distance from a point to a line segment as described
// here:
// http://geomalgorithms.com/a02-_lines.html#Distance%20to%20Ray%20or%20Segment

func distToEdge(k, i, j Node) float64 {
	//k is node, i to j is an edge
	u := Node{0, j.X - i.X, j.Y - i.Y}
	v := Node{0, k.X - i.X, k.Y - i.Y}

	c1 := dot(v, u)
	if c1 <= 0 {
		return distance(k, i)
	}
	c2 := dot(u, u)
	if c2 <= c1 {
		return distance(k, j)
	}
	b := c1 / c2
	pt := Node{0, i.X + b*u.X, i.Y + b*u.Y}
	return distance(k, pt)
}

//-- helper functions, may not be needed.
func dot(u, v Node) float64 {
	return u.X*v.X + u.Y*v.Y
}

func norm(n Node) float64 {
	return math.Sqrt(dot(n, n))
}
