package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fn := "Random100.tsp"
	// fn := "Random10.tsp"
	if len(os.Args) > 1 {
		fn = os.Args[1]
	}
	initPopSize := 100
	// initPopSize := 5
	if len(os.Args) > 2 {
		initPopSize, _ = strconv.Atoi(os.Args[2])
	}

	generations := 300
	if len(os.Args) > 3 {
		generations, _ = strconv.Atoi(os.Args[3])
	}

	elitePerc := float64(.2)
	if len(os.Args) > 4 {
		elitePerc, _ = strconv.ParseFloat(os.Args[4], 64)
		// eliteSize = int(perc * float64(initPopSize))
	}
	//FIXME:

	//Default: 5% mutation rate.
	var mutationRate float64 = 0.5
	if len(os.Args) > 5 {
		mutationRate, _ = strconv.ParseFloat(os.Args[5], 64)
	}

	fmt.Printf("%s, popsize: %d, gens: %v, eliteperc: %v, mutation: %v\n", fn, initPopSize, generations, elitePerc, mutationRate)

	nodes := loadNodes(fn)
	// _, gens := geneticAlgo(nodes, initPopSize, eliteSize, mutationRate, generations)

	// writeGraphToFile(&gens)

	eliteSize := int(elitePerc * float64(len(nodes)))
	if eliteSize > len(nodes) {
		eliteSize = len(nodes)
	}

	// // init pop
	// initpop := initPop(initPopSize, nodes)
	// // sort init pop
	// sorted := sortByFitness(initpop)
	// fmt.Printf("sorted:\n[")
	// for _, v := range sorted {
	// 	fmt.Printf("%v ", v.tripID)
	// }
	// fmt.Printf("\b]\n")
	// // roulette select from pop
	// selected := selectPop(sorted, eliteSize)
	// fmt.Println("selected:")
	// fmt.Println(selected)
	// // create mating pool from selected
	// pool := matingPool(initpop, selected)
	// for i, v := range pool {
	// 	fmt.Println(i, v)
	// }
	// fmt.Printf("\n")
	// //mutate
	// mutatePop(pool, mutationRate)
	// fmt.Println("mutated:")
	// for i, v := range pool {
	// 	fmt.Println(i, v)
	// }
	// fmt.Printf("\n")

	// // test breed first two from matking pool
	// fmt.Println("p1:", pool[0])
	// fmt.Println("p2:", pool[1])
	// child := breed(pool[0], pool[1])
	// fmt.Println("ch:", child)
	// children := breedPop(pool, eliteSize)
	// mutatePop(children, mutationRate)
	// fmt.Println("gen1:")
	// for _, v := range children {
	// 	fmt.Println(v)
	// }

	// nextgen := nextGen(children, eliteSize, mutationRate)
	// fmt.Println("gen2")
	// for _, v := range nextgen {
	// 	fmt.Println(v)
	// }

	// fmt.Println("putting it together:")
	start := time.Now()
	var best *Trip
	var gens []Generation
	// for {
	// 	best, gens = geneticAlgo(nodes, initPopSize, elitePerc, mutationRate, generations)
	// 	fmt.Println(best.cost)
	// 	fmt.Println(len(gens))
	// 	drawP(*best.path, nodes)
	// 	if best.cost <= 1000.0 {
	// 		drawPath(best, nodes, "path100")
	// 		break
	// 	}
	// }
	best, gens = geneticAlgo(nodes, initPopSize, elitePerc, mutationRate, generations)
	fmt.Println("runtime:", time.Since(start))
	fmt.Println(best.cost)
	fmt.Println(len(gens))
	// drawPath(best, nodes, "path100")
	writeGraphToFile(&gens)

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
