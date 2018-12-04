package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

// Start of main program here. Calls other functions to obtain optimal trip
// in O(n!) time. brute force approach.
func bruteforce(filename string) *Trip {
	//read file
	nodes := loadNodes(filename)
	fmt.Println("size: ", len(nodes))
	// pairs := combinations(nodes)
	// fmt.Println(pairs)

	// get all permutations
	trip := Permutate(nodes)
	// fmt.Println(trip)
	return trip
}

//Passes filename, opens file, reads file and stores them in a map
//with the city index as the index in the map/slice pointing to a Node which
//consists of the city index as well as the x & y coords.
// nodes := make(map[int]Node)
func loadNodes(filename string) []Node {

	var nodes []Node
	//open file, handle error
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() //closes file at end of function

	scanner := bufio.NewScanner(file)

	//iterate thru file reading line by line, discarding everything before and
	//include the NODE_COORD_SECTION line.
	for scanner.Scan() {
		if scanner.Text() == "NODE_COORD_SECTION" {
			break
		}
	}

	//the rest of the lines in the file will bee coordinates, read lines
	//separate by spaces, load onto Node struct and add to map.
	for scanner.Scan() {
		var idx int
		var x, y float64

		_, err := fmt.Sscanf(scanner.Text(), "%d %f %f", &idx, &x, &y)
		if err != nil {
			log.Fatal(err)
		}

		//add Node to nodes map or slice.
		// nodes[idx] = Node{idx, x, y}
		curr := Node{idx, x, y}
		nodes = append(nodes, curr)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nodes
}

// calculates the distance between two nodes.
func distance(a, b Node) float64 {
	//hard way, see if the math package has a way of doing this without me
	//implementing it. dist = sqrt((b.x - a.x)^2 + (b.y - a.y^2))
	dist := math.Hypot(b.X-a.X, b.Y-a.Y)
	return dist
}

// Permutate finds all current permutations from a given array. Uses helper
// function permutate() below to obtain all permutations using recursion.
func Permutate(n []Node) *Trip {
	shortest := &Trip{&[]Node{}, math.MaxFloat64}
	pairs := combinations(n)
	permutate(n, len(n), shortest, pairs)
	return shortest
}

// permutation utilizes "Heaps Algorithm" to find all possible permutations
// given a slice of Nodes.  more info at:
// https://en.wikipedia.org/wiki/Heap%27s_algorithm
// recursive permutate function: (not good with concurrency, not finished)
func permutate(n []Node, size int, shortest *Trip, p []Pair) {
	if size == 1 {
		//when size == 1 a permutation is found, calculate distance and compare
		//to current shortest trip, if the current permutation total travel cost
		//is shorter than current shortest, update current shortest to current
		//permutation.
		var dist float64
		//THIS PART OF CODE IS ASSUMNING ALL TRIPS START WITH FIRST CITY,
		//COMMENT OUT AND UNCOMMEND CODE BELOW TO FIND TRUE SHORTEST TRIP
		if n[0].ID == 1 {
			if len(n) == 2 {
				dist = distance(n[0], n[1])
			} else {
				for i := 1; i < len(n); i++ {
					dist += distance(n[i-1], n[i])
				}
			}
			dist += distance(n[len(n)-1], n[0])
			if dist < shortest.cost {
				path := make([]Node, len(n))
				copy(path, n)
				*shortest = Trip{&path, dist}
			}
		}
		//UNCOMMENT THIS CODE AND COMMENT ABOVE IF STATEMENT IF WANTING TO FIND
		//TRUE SHORTEST TRIP.
		// if len(n) == 2 {
		//      dist = distance(n[0], n[1])
		//  } else {
		//      for i := 1; i < len(n); i++ {
		//          dist += distance(n[i-1], n[i])
		//      }
		//  }
		//  dist += distance(n[len(n)-1], n[0])
		//  if dist < shortest.cost {
		//      path := make([]Node, len(n))
		//      copy(path, n)
		//      *shortest = Trip{&path, dist}
		//  }
	} else {
		for i := 0; i < size; i++ {
			permutate(n, size-1, shortest, p)
			if size%2 == 0 {
				n[i], n[size-1] = n[size-1], n[i]
			} else {
				n[0], n[size-1] = n[size-1], n[0]
			}
		}
	}
}

//combination formula nCr that works with this project, not being used at the
//moment, will be used for optimization in future.
func combinations(n []Node) []Pair {
	var pairs []Pair
	for i := 0; i < len(n)-1; i++ {
		for j := i + 1; j < len(n); j++ {
			pair := Pair{n[i], n[j], distance(n[i], n[j])}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

//not being used at the moment.
func getDist(a, b Node) float64 {
	pair := Pair{a, b, 0.0}
	pair.sort()
	pair.distance()
	return pair.dist
}
