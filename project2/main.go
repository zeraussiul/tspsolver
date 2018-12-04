package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

func main() {

	nodes := loadNodes("11PointDFSBFS.tsp")
	dgraph := buildDGraph(nodes)

	start := time.Now()
	bfspath := bfs(dgraph, 11)
	bfsDistance := pathDistance(bfspath)
	end := time.Since(start)
	fmt.Println("bfs:", bfspath, "distance:", bfsDistance)
	fmt.Println("running time:", end)

	start = time.Now()
	dfspath := dfs(dgraph, 11)
	dfsDistance := pathDistance(dfspath)
	end = time.Since(start)
	fmt.Println("dfs:", dfspath, "distance:", dfsDistance)
	fmt.Println("running time:", end)

	// make trip strucs to pass to drawPath() function that draws the trip as a
	// png file
	bfsTrip := &Trip{&bfspath, bfsDistance}
	dfsTrip := &Trip{&dfspath, dfsDistance}

	// draw the trips into png files
	drawPath(bfsTrip, nodes, "bfstrip")
	drawPath(dfsTrip, nodes, "dfstrip")

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

// calculates the distance between two nodes.
func distance(a, b Node) float64 {
	//hard way, see if the math package has a way of doing this without me
	//implementing it. dist = sqrt((b.x - a.x)^2 + (b.y - a.y^2))
	dist := math.Hypot(b.X-a.X, b.Y-a.Y)
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
