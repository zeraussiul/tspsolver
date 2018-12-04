package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func woc(nodes []Node, crowdSize int) *Trip {
	//crowd holds the trips that will be used to aggregate.
	crowd := make(chan *Trip, crowdSize)

	//done channel blocks main goroutine from exiting before other goroutines
	//are done.
	done := make(chan bool, crowdSize)
	for i := 0; i < crowdSize; i++ {
		go proc_input(nodes, done, crowd)
	}
	for i := 0; i < crowdSize; i++ {
		<-done
	}
	trips := []*Trip{}
	for i := 0; i < crowdSize; i++ {
		trips = append(trips, <-crowd)
	}
	close(crowd)

	return aggregate(nodes, trips)
}

func proc_input(nodes []Node, done chan<- bool, crowd chan<- *Trip) {

	initPopSize := 150
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

	//Default: 50% mutation rate.
	var mutationRate float64 = 0.1
	if len(os.Args) > 5 {
		mutationRate, _ = strconv.ParseFloat(os.Args[5], 64)
	}

	eliteSize := int(elitePerc * float64(len(nodes)))
	if eliteSize > len(nodes) {
		eliteSize = len(nodes)
	}

	// fmt.Println("pop:", initPopSize, "gens:", generations, "eliteperc:", elitePerc, "mutation:", mutationRate)
	// start := time.Now()
	var best *Trip
	// var gens []Generation
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
	best = geneticAlgo(nodes, initPopSize, elitePerc, mutationRate, generations)
	// fmt.Println("runtime:", time.Since(start))
	// fmt.Println(best.cost)
	// fmt.Println(len(gens))
	// drawPath(best, nodes, "path100")
	// writeGraphToFile(&gens)
	crowd <- best
	done <- true
}

func aggregate(nn []Node, crowd []*Trip) *Trip {
	drawP(*crowd[0].path, nn, "testCrowd")
	fmt.Println("GA  Cost:", crowd[0].cost)

	crowdMatrix := createMatrix(crowd)
	edgeList := createEdgeList(crowdMatrix)
	// fmt.Println("Edge List:", edgeList)
	subpaths := createSubPaths(nn, edgeList)
	// fmt.Println("Subpaths:", subpaths)

	//we will add subpaths to trip connecting them via greedy algo.
	trip := []Node{}
	currSub := subpaths[0]
	trip = append(trip, currSub...)
	//remove subpath added to tree
	subpaths = append(subpaths[:0], subpaths[1:]...)

	for len(subpaths) != 0 {
		var currIdx, nextIdx, subpathIdx int
		shortest := math.MaxFloat64
		for k, s := range subpaths {
			for i := 0; i < len(trip); i += len(trip) - 1 {
				for j := 0; j < len(s); j += len(s) - 1 {
					dist := distance(trip[i], s[j])
					if dist < shortest {
						shortest = dist
						currIdx = i
						nextIdx = j
						subpathIdx = k
					}
				}
			}
		}
		subpath := subpaths[subpathIdx]
		//checks where to append the chosen subpath, and whether reversing the
		//subpath is needed first.
		if currIdx == 0 && nextIdx != 0 {
			//NOTE: append subpath to front without reversing
			trip = append(subpath, trip...)
		} else if currIdx != 0 && nextIdx == 0 {
			//NOTE:	append supath to back without reversing
			trip = append(trip, subpath...)
		} else if currIdx == 0 && nextIdx == 0 {
			//NOTE:	reverse subpath and append to front
			reverse(subpath)
			trip = append(subpath, trip...)
		} else {
			//NOTE: reverse subpath and append to back
			reverse(subpath)
			trip = append(trip, subpath...)
		}

		//remove path added to trip from subpath
		subpaths = append(subpaths[:subpathIdx], subpaths[subpathIdx+1:]...)

	}

	shortest := math.MaxFloat64
	var rotateAt int
	for i, v := range trip {
		if i == len(trip)-1 {
			break
		}
		dist := distance(v, trip[i+1])
		if dist < shortest {
			shortest = dist
			rotateAt = i + 1
		}

	}

	rotate(trip, len(trip)-rotateAt)
	city := Node{trip[0].ID, trip[0].X, trip[0].Y}
	trip = append(trip, city)

	final := &Trip{&trip, pathDistance(trip)}

	return final
}

func createMatrix(crowd []*Trip) [][]int {
	n := len(*crowd[0].path)

	//https://stackoverflow.com/questions/39804861/what-is-a-concise-way-to-create-a-2d-slice-in-go
	//crowdMatrix will hold the edges the crowd decided on, from here we will increment +1 for every edge found on each crowd until finished.
	crowdMatrix := make([][]int, n)
	for i := range crowdMatrix {
		crowdMatrix[i] = make([]int, n)
	}

	//FIXME: iterate thru each trip inside crowd, populating the crowdMatrix, not efficient, will look at other solutions later.
	for _, c := range crowd {
		trip := *c.path
		for k := 0; k < n-1; k++ {
			i, j := trip[k].ID, trip[k+1].ID
			if j < i {
				crowdMatrix[j][i]++
			} else {
				crowdMatrix[i][j]++
			}
		}
	}

	//print matrix, fix spacing later.
	// for i, v := range crowdMatrix {
	// 	if i == 0 {
	// 		//discards first array, not needed
	// 		continue
	// 	}
	// 	fmt.Printf("%2d %v\n", i, v)
	// }

	return crowdMatrix
}

func createEdgeList(crowdMatrix [][]int) []int {
	n := len(crowdMatrix)
	//edgeList will hold the most agreed edge according to crowd. the index of
	//edgeList is the current city, and the value it holds is the city it
	//connects to.
	edgeList := make([]int, n)
	//addedCities keeps track of cities already added to edgeList so two cities
	//arent added.
	addedCities := make(map[int]bool)
	for i := range edgeList {
		if i == 0 {
			continue
		}
		currentMax := 0
		//iterate thru crowdMatrix
		for j := 1; j < n; j++ {
			if currentMax < crowdMatrix[i][j] {

				if addedCities[j] {
					continue
				}
				currentMax = crowdMatrix[i][j]
				edgeList[i] = j
			}
		}
		addedCities[edgeList[i]] = true

	}
	return edgeList
}

func createSubPaths(nn []Node, edgeList []int) [][]Node {
	subpaths := [][]Node{}
	addedToSub := make(map[int]bool)

	for i, v := range edgeList {
		if v == 0 || addedToSub[i] {
			continue
		}
		sub := []Node{}
		sub = append(sub, nn[i-1])
		addedToSub[i] = true
		var nextIdx int
		for v != 0 {
			nextIdx = v
			sub = append(sub, nn[nextIdx-1])
			addedToSub[nextIdx] = true
			v = edgeList[nextIdx]
		}

		subpaths = append(subpaths, sub)
	}
	return subpaths
}

func rotate(nn []Node, k int) {
	if len(nn) <= 1 {
		return
	}
	//do three total rotations, one from len(nums) - k to end
	//second from 0 to len(nums) - 1 - l
	//finally rotate all
	k = k % len(nn) //in case the array is smaller than k.
	size := len(nn)
	reverseArray(nn, size-k, size-1)
	reverseArray(nn, 0, size-1-k)
	reverseArray(nn, 0, size-1)

}

//helper func to rotate()
func reverseArray(nn []Node, start, end int) {
	//reverses array in place
	//if length of array is less than or equal to one we assume
	//there is no need to reverse the array, return.
	if len(nn) <= 1 {
		return
	}
	for start < end {
		nn[start], nn[end] = nn[end], nn[start]
		start++
		end--
	}
}
