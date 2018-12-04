package main

import (
	"math"
)

//DFS implementation from:
// https://www.hackerearth.com/practice/algorithms/graphs/depth-first-search/tutorial/
//and: https://en.wikipedia.org/wiki/Depth-first_search

func dfs(dd []*DGraph, target int) []Node {
	//Reset graph
	for _, v := range dd {
		v.visited = WHITE
		v.dist = math.MaxInt64
		v.prev = nil
	}

	//init stack
	s := make(stack, 0)
	//push first element onto stack
	s = s.Push(dd[0])

	var curr *DGraph
	//iterate thru stack til empty
	for s.Len() != 0 {
		// fmt.Println(s)
		s, curr = s.Pop()
		curr.visited = BLACK

		//push all adjacent cities to current onto stack
		for _, v := range curr.to {
			if v.visited == WHITE {
				v.visited = BLACK
				v.prev = curr
				s = s.Push(v)
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

	reverse(path)

	return path
}
