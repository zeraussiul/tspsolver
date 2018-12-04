package main

// Directed Graph built as follows:
// 1 to 2, 3, 4
// 2 to 3
// 3 to 4, 5
// 4 to 5, 6, 7
// 5 to 7, 8
// 6 to 8
// 7 to 9, 10
// 8 to 9, 10, 11
// 9 to 11
// 10 to 11
//read from 11PointDFSBFS.tsp provided with the project.

func buildDGraph(n []Node) []*DGraph {

	// dgraph := []DGraph{
	// 	DGraph{n[0], []*Node{&n[1], &n[2], &n[3]}},
	// 	DGraph{n[1], []*Node{&n[2]}},
	// 	DGraph{n[2], []*Node{&n[3], &n[4]}},
	// 	DGraph{n[3], []*Node{&n[4], &n[5], &n[6]}},
	// 	DGraph{n[4], []*Node{&n[6], &n[7]}},
	// 	DGraph{n[5], []*Node{&n[7]}},
	// 	DGraph{n[6], []*Node{&n[8], &n[9]}},
	// 	DGraph{n[7], []*Node{&n[8], &n[9], &n[10]}},
	// 	DGraph{n[8], []*Node{&n[10]}},
	// 	DGraph{n[9], []*Node{&n[10]}},
	// }
	v11 := DGraph{n[10], WHITE, 0, nil, nil}
	v10 := DGraph{n[9], WHITE, 0, nil, []*DGraph{&v11}}
	v09 := DGraph{n[8], WHITE, 0, nil, []*DGraph{&v11}}
	v08 := DGraph{n[7], WHITE, 0, nil, []*DGraph{&v09, &v10, &v11}}
	v07 := DGraph{n[6], WHITE, 0, nil, []*DGraph{&v09, &v10}}
	v06 := DGraph{n[5], WHITE, 0, nil, []*DGraph{&v08}}
	v05 := DGraph{n[4], WHITE, 0, nil, []*DGraph{&v07, &v08}}
	v04 := DGraph{n[3], WHITE, 0, nil, []*DGraph{&v05, &v06, &v07}}
	v03 := DGraph{n[2], WHITE, 0, nil, []*DGraph{&v04, &v05}}
	v02 := DGraph{n[1], WHITE, 0, nil, []*DGraph{&v03}}
	v01 := DGraph{n[0], WHITE, 0, nil, []*DGraph{&v02, &v03, &v04}}

	dgraph := []*DGraph{&v01, &v02, &v03, &v04, &v05, &v06, &v07, &v08, &v09, &v10, &v11}

	return dgraph
}
