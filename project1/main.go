package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	//The program is capable of receiving multiple filenames and
	//running them, printing the est. running time to console.
	var shortestTrip *Trip
	start := time.Now()
	if len(os.Args) == 2 {

		shortestTrip = bruteforce(os.Args[1])
		fmt.Println("running time:", time.Since(start))
	} else {
		for i := 1; i < len(os.Args); i++ {

			shortestTrip = bruteforce(os.Args[i])
			fmt.Println("running time:", time.Since(start))
		}
	}
	fmt.Println(shortestTrip)
	// shortest trip found, draw to PNG file in a graph
	// pass a *Trip as well as a filename to save to. in our case we will
	// pass the arg passed onto main and append .PNG to it, ex: Random4.tsp.png
	drawPath(shortestTrip, os.Args[1])
}
