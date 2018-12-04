package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const ext = ".tsp"

func main() {
	rand.Seed(time.Now().UnixNano())
	fn := "Random11.tsp"
	// fn := "Random10.tsp"
	if len(os.Args) > 1 {
		fn = os.Args[1]
		fe := filepath.Ext(fn)
		if fe != ext {
			log.Fatalf("Wrong filetype. Expected: %v, got: %v", ext, fe)
		}
	}
	nodes := loadNodes(fn)
	//default crowd size = 5
	crowdSize := 5
	if len(os.Args) > 6 {
		c, err := strconv.Atoi(os.Args[6])
		if err != nil {
			log.Fatalf("%v", err)
		}
		crowdSize = c
	}

	//main entrance to Wisdom of Crowds problem
	start := time.Now()
	fmt.Println("Crowd Size:", crowdSize)
	best := woc(nodes, crowdSize)
	fmt.Println("WOC Cost:", best.cost)
	fmt.Println("Running time:", time.Since(start))
	drawP(*best.path, nodes, "testFinal")

	// fmt.Println(best)
}

//-------------HELPER FUNCTIONS----------------------------//

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

// calculates the distance between two nodes, or an Edge.
func distance(a, b Node) float64 {
	//hard way, see if the math package has a way of doing this without me
	//implementing it. dist = sqrt((b.x - a.x)^2 + (b.y - a.y^2))
	dist := math.Hypot(b.X-a.X, b.Y-a.Y)
	// dist := math.Hypot(a.X-b.X, a.Y-b.Y)
	return dist
}

// reverses a slice of nodes.
func reverse(nn []Node) {
	for i := len(nn)/2 - 1; i >= 0; i-- {
		opp := len(nn) - 1 - i
		nn[i], nn[opp] = nn[opp], nn[i]
	}
}

// returns the total distance traveled on the trip passed as parameter.
func pathDistance(nn []Node) float64 {
	var dist float64
	for i := 1; i < len(nn); i++ {
		dist += distance(nn[i-1], nn[i])
	}

	return dist
}
