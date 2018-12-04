package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

//----Create Population----
//creates a random trip, relies on math/rand for the shuffle func
func createTrip(nn []Node) *Trip {

	//copy original city list onto Node slice called path
	path := make([]Node, len(nn))
	copy(path, nn)
	//shuffle path and add starting city to end of path
	rand.Shuffle(len(path), func(i, j int) {
		path[i], path[j] = path[j], path[i]
	})
	start := Node{path[0].ID, path[0].X, path[0].Y}
	path = append(path, start)

	trip := &Trip{&path, pathDistance(path)}

	return trip
}

//initPop creates a population of size containing random trips
func initPop(size int, nn []Node) []*Trip {
	//seed rand for createTrip() to work properly
	// rand.Seed(time.Now().UnixNano())
	pop := []*Trip{}
	for i := 0; i < size; i++ {
		pop = append(pop, createTrip(nn))
	}
	return pop
}

//fitness is inverse of cost of trip
func fitness(t *Trip) float64 {
	return 1 / t.cost
}

//---Determine Fitness----
//sortByFitness returns a list of trips sorted by their fitness.
// func sortByFitness(pop []*Trip) []Fitness {
// 	fitnessResults := []Fitness{}
// 	for _, v := range pop {
// 		fitnessResults = append(fitnessResults, Fitness{v, fitness(v)})
// 	}
// 	sort.Sort(ByFitness(fitnessResults))
// 	return fitnessResults
// }
func sortByFitness(pop []*Trip) []Fitness {
	fitnessResults := []Fitness{}
	for i, v := range pop {
		fitnessResults = append(fitnessResults, Fitness{i, fitness(v)})
	}
	sort.Sort(ByFitness(fitnessResults))
	return fitnessResults
}

//---Selecting mating pool---
func selectPop(ff []Fitness, eliteSize int) []int {
	// return rouletteWheelSelection(ff, eliteSize)
	return rankedSelection(ff, eliteSize)
}

//Roulette Wheel Selection with eiltism
//1.[Sum]: of fitness in pop, S
//2.[Select]: generate rand number from 0 to S, R
//3.[Loop]: go thru pop and sum fitnesses from 0 to S, when sum is greater than //R, stop and return chromosome of where we are.
func rouletteWheelSelection(ff []Fitness, eliteSize int) []int {
	selectedPopID := []int{}

	s := fitnessSum(&ff)
	//first, add the best pops from 0 to elitesize
	for i := 0; i < eliteSize; i++ {
		selectedPopID = append(selectedPopID, ff[i].tripID)
	}
	//then randomly pick the rest up to popsize utilizing roulette wheel selection
	for i := 0; i < len(ff)-eliteSize; i++ {
		//select random number between o and s
		r := randRange(0.0, s)
		//loop thru pop and sum from 0 to S until partial sum > than r
		var partial float64
		for _, f := range ff {
			partial += f.fitness
			if partial > r {
				selectedPopID = append(selectedPopID, f.tripID)
				break
			}
		}
	}
	return selectedPopID
}

//https://cs.stackexchange.com/questions/89886/how-is-rank-selection-better-than-random-selection-and-rws
func rankedSelection(ff []Fitness, eliteSize int) []int {
	selectedPopID := []int{}

	// s := fitnessSum(&ff)
	//first, add the best pops from 0 to elitesize
	for i := 0; i < eliteSize; i++ {
		selectedPopID = append(selectedPopID, ff[i].tripID)
	}

	bias := 1.5
	for i := 0; i < len(ff)-eliteSize; i++ {
		idx := int(float64(len(ff)) * (bias - math.Sqrt(bias*bias-4.0*(bias-1)*rand.Float64())) / 2.0 / (bias - 1))
		selectedPopID = append(selectedPopID, ff[idx].tripID)
	}

	return selectedPopID
}

func fitnessSum(ff *[]Fitness) float64 {
	var sum float64
	for _, f := range *(ff) {
		sum += f.fitness
	}
	return sum
}

func randRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

//matingPool takes in the population as well as the int slice of selectedPopIds and creates a mating pool
//MAY NEED WORK!
func matingPool(pop []*Trip, selectedPopID []int) []*Trip {
	matingpool := []*Trip{}

	for _, v := range selectedPopID {
		matingpool = append(matingpool, pop[v])
	}

	return matingpool
}

//mating takes place by picking a random gene from one parent (i, j) adding it
//to child trip, and adding missing cities using ordered crossover.
func breed(p1, p2 *Trip) *Trip {

	childTrip := make([]Node, len(*p1.path))
	p1path := *p1.path
	p2path := *p2.path

	geneA := rand.Intn(len(*p1.path))
	geneB := rand.Intn(len(*p1.path))
	var start, end int
	if geneA > geneB {
		geneA, geneB = geneB, geneA
	}
	start, end = geneA, geneB

	inChild := make(map[int]int)

	//add onto chilTrip the genes from p1 from start to end.
	for i := start; i < end; i++ {
		childTrip[i] = p1path[i]
		inChild[p1path[i].ID] = i
	}
	//add remaining genes from p2, while not duplicating cities already added from p1
	idx := 0
	for i := 0; i < len(p2path)-1; i++ {
		for childTrip[idx].ID != 0 {
			idx++
		}
		if _, ok := inChild[p2path[i].ID]; !ok {
			childTrip[idx] = p2path[i]
			idx++
		}
	}

	//build a copy of the first city visited and complete the trip by adding it
	//to the last element of the childTrip slice.s
	lasCity := Node{childTrip[0].ID, childTrip[0].X, childTrip[0].Y}
	// childTrip = append(childTrip, lasCity)
	childTrip[len(childTrip)-1] = lasCity

	child := &Trip{&childTrip, pathDistance(childTrip)}

	return child
}

func breedPop(matingpool []*Trip, eliteSize int) []*Trip {
	children := []*Trip{}
	for i := 0; i < eliteSize; i++ {
		children = append(children, matingpool[i])
	}

	length := len(matingpool) - eliteSize

	for i := 0; i < length; i++ {
		p1 := rand.Intn(len(matingpool) - 1)
		p2 := rand.Intn(len(matingpool) - 1)
		// if p1 > p2 {
		// 	p1, p2 = p2, p1
		// }
		child := breed(matingpool[p1], matingpool[p2])
		children = append(children, child)
	}

	return children
}

//mutate by swapping two random cities from the trip, edit individual trip by
//reference.
//NOTE: only mutates if it improves overall trip, otherwise it does not
func mutate(trip **Trip, mutationRate float64) {
	path := *(*trip).path
	for i := range path {
		//skip start and end city
		if i == 0 || i == len(path)-1 {
			continue
		}
		r := rand.Float64()
		if r < mutationRate {
			// get index thats not 0 or last index in array, we dont want to randomly swap the start and end city.
			j := rand.Intn(len(path)-2) + 1
			path[i], path[j] = path[j], path[i]
			//checks if new path cost is better than old one, if its not, dont mutate. not efficient rn
			newcost := pathDistance(path)
			if newcost < (*trip).cost {
				(*trip).cost = newcost
				// fmt.Println("mutated, new val:", (*trip).cost)
			} else {
				path[i], path[j] = path[j], path[i]
			}
		}
	}
}

//run mutation algo on entire population
func mutatePop(pop []*Trip, mutationRate float64) {
	for _, p := range pop {
		mutate(&p, mutationRate)
	}
}

func nextGen(currentpop []*Trip, eliteSize int, mutationRate float64) []*Trip {
	sortedByFitness := sortByFitness(currentpop)
	selectedPopID := selectPop(sortedByFitness, eliteSize)
	matingpool := matingPool(currentpop, selectedPopID)
	children := breedPop(matingpool, eliteSize)
	//mutate children
	mutatePop(children, mutationRate)
	return children
}

func geneticAlgo(nn []Node, popsize int, elitePerc float64, mutationRate float64, generations int) (*Trip, []Generation) {
	//seed rand
	// rand.Seed(time.Now().UnixNano())

	//run algo:
	eliteSize := int(elitePerc * float64(len(nn)))
	if eliteSize > len(nn) {
		eliteSize = len(nn)
	}

	if mutationRate > 1.0 {
		mutationRate = 1.0
		fmt.Println("mutationRate changed to", mutationRate)
	}

	// var bestTrip *Trip
	gens := []Generation{}
	pop := initPop(popsize, nn)
	sortedByFitness := sortByFitness(pop)
	first := pop[sortedByFitness[0].tripID].cost
	last := pop[sortedByFitness[len(sortedByFitness)-1].tripID].cost
	avg := (first + last) / 2
	gens = append(gens, Generation{0, first, avg, last})

	for i := 1; i <= generations; i++ {
		pop = nextGen(pop, eliteSize, mutationRate)
		sortedByFitness = sortByFitness(pop)
		first = pop[sortedByFitness[0].tripID].cost
		last = pop[sortedByFitness[len(sortedByFitness)-1].tripID].cost
		avg = (first + last) / 2
		gens = append(gens, Generation{i, first, avg, last})
	}
	best := pop[0]
	return best, gens
}
