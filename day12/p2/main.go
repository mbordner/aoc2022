package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"github.com/mbordner/aoc2022/common/geom"
	"github.com/mbordner/aoc2022/common/graph"
	"github.com/mbordner/aoc2022/common/graph/djikstra"
)

var (
	bb         = &geom.BoundingBox{}
	nodePosMap = make(map[geom.Pos]geom.Pos)
)

func main() {
	G, S, E, lowNodes := getMapData("../data.txt")

	shortestDistance := G.GetNodeCount()
	var bestStart *graph.Node

	for _, ln := range lowNodes {
		shortestPaths := djikstra.GenerateShortestPaths(G, ln)
		shortestPath, _ := shortestPaths.GetShortestPath(E)

		dis := len(shortestPath)
		if dis == 0 {
			continue
		}

		if dis < shortestDistance {
			shortestDistance = dis
			bestStart = ln
		}
	}

	fmt.Println("S is at: ", S.GetID())
	fmt.Println("E is at: ", E.GetID())

	if bestStart == S {
		fmt.Println("best start is at S")
	} else {
		fmt.Println("best start is at  ", bestStart.GetID())
	}

	fmt.Println("Shortest distance: ", shortestDistance)
}

func getElevation(char rune) int {
	if char == 'S' {
		char = 'a'
	} else if char == 'E' {
		char = 'z'
	}
	return int(char - 'a')
}

func getMapData(path string) (*graph.Graph, *graph.Node, *graph.Node, []*graph.Node) {

	g := graph.NewGraph()

	var S, E *graph.Node
	lowNodes := make([]*graph.Node, 0, 10)

	lines, _ := file.GetLines(path)
	for j, line := range lines {
		for i, char := range line {

			pos := geom.Pos{X: i, Y: j, Z: getElevation(char)}
			nodePosMap[geom.Pos{X: i, Y: j}] = pos

			bb.Extend(pos)
			n := g.CreateNode(pos)
			if char == 'S' {
				S = n
			} else if char == 'E' {
				E = n
			}

			if pos.Z == 0 {
				lowNodes = append(lowNodes, n)
			}

		}
	}

	for j := 0; j <= bb.MaxY; j++ {
		for i := 0; i <= bb.MaxX; i++ {
			nPos := nodePosMap[geom.Pos{X: i, Y: j}]
			n := g.GetNode(nPos)

			// north
			if oPos, exists := nodePosMap[geom.Pos{X: i, Y: j - 1}]; exists {
				o := g.GetNode(oPos)
				if o != nil {
					if oPos.Z-nPos.Z <= 1 {
						e := n.AddEdge(o, 1)
						e.AddProperty("dir", geom.North)
					}
				}
			}

			// east
			if oPos, exists := nodePosMap[geom.Pos{X: i + 1, Y: j}]; exists {
				o := g.GetNode(oPos)
				if o != nil {
					if oPos.Z-nPos.Z <= 1 {
						e := n.AddEdge(o, 1)
						e.AddProperty("dir", geom.East)
					}
				}
			}

			// south
			if oPos, exists := nodePosMap[geom.Pos{X: i, Y: j + 1}]; exists {
				o := g.GetNode(oPos)
				if o != nil {
					if oPos.Z-nPos.Z <= 1 {
						e := n.AddEdge(o, 1)
						e.AddProperty("dir", geom.South)
					}
				}
			}

			// west
			if oPos, exists := nodePosMap[geom.Pos{X: i - 1, Y: j}]; exists {
				o := g.GetNode(oPos)
				if o != nil {
					if oPos.Z-nPos.Z <= 1 {
						e := n.AddEdge(o, 1)
						e.AddProperty("dir", geom.West)
					}
				}
			}
		}
	}

	return g, S, E, lowNodes

}
