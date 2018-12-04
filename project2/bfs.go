package main

import "math"

//BFS implementation derived from pseudo code: https://www.hackerearth.com/practice/algorithms/graphs/breadth-first-search/tutorial/

func bfs(dd []*DGraph, target int) []Node {
	//Reset graph
	for _, v := range dd {
		v.visited = WHITE
		v.dist = math.MaxInt64
		v.prev = nil
	}

	//initiate queue
	q := make(stack, 0)

	//push source Node onto queue
	dd[0].visited = BLACK
	q = q.Queue(dd[0])
	var curr *DGraph

	//iterate thru queue until empty, which means we've found all paths
	//originating from source
	for q.Len() != 0 {
		//dequeue next element
		q, curr = q.Dequeue()

		//add adjacenent cities to queue
		for _, v := range curr.to {
			if v.visited == WHITE {
				v.visited = BLACK
				v.prev = curr
				q = q.Push(v)
			}
		}
	}

	//work backwards to trace path back to source.
	path := make([]Node, 0)

	r := dd[len(dd)-1]
	for r != nil {
		path = append(path, r.node)
		r = r.prev
	}

	//because we worked "backwards" when creating the array holding the path in
	//the above section of the code, we have to reverse the array to make the
	//trip in the right order.""
	reverse(path)

	return path
}
