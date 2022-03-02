package waterjug

import (
	"fmt"
)

type Operation string

const (
	OpFillX      = "Fill bucket X"
	OpFillY      = "Fill bucket Y"
	OpDumpX      = "Dump bucket X"
	OpDumpY      = "Dump bucket Y"
	OpTransferXY = "Transfer bucket X to bucket Y"
	OpTransferYX = "Transfer bucket Y to bucket X"
)

type Node struct {
	X  int
	Y  int
	Op string
}

func SearchSolution(x, y, z int) ([]Node, bool) {
	initialNode := Node{0, 0, ""}
	paths := [][]Node{{initialNode}}
	visitedNodes := map[string]Node{}
	for len(paths) > 0 {
		path := paths[0]
		paths = paths[1:]
		lastNode := path[len(path)-1]
		visitedNodes[getIndex(lastNode)] = lastNode
		if isSolution(path, z) {
			return path, true
		}
		nextMoves := nextTransitions(x, y, path, visitedNodes)
		for i := range nextMoves {
			paths = append(paths, nextMoves[i])
		}
	}
	return nil, false
}

func getIndex(n Node) string {
	return fmt.Sprint(n.X, ":", n.Y)
}

func isSolution(path []Node, z int) bool {
	if len(path) == 0 {
		return false
	}
	return path[len(path)-1].X == z || path[len(path)-1].Y == z
}

func beenThere(node Node, visitedNodes map[string]Node) bool {
	_, ok := visitedNodes[getIndex(node)]
	return ok
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// find next possible moves, except visited nodes.
func nextTransitions(x, y int, path []Node, visitedNodes map[string]Node) [][]Node {
	var (
		result    [][]Node
		nextNodes []Node
		aMax      = x
		bMax      = y
		a         = path[len(path)-1].X
		b         = path[len(path)-1].Y
	)

	node := Node{aMax, b, OpFillX}
	if !beenThere(node, visitedNodes) {
		nextNodes = append(nextNodes, node)
	}
	node = Node{a, bMax, OpFillY}
	if !beenThere(node, visitedNodes) {
		nextNodes = append(nextNodes, node)
	}
	node = Node{intMin(aMax, a+b), b - (intMin(aMax, a+b) - a), OpTransferYX}
	if !beenThere(node, visitedNodes) {
		nextNodes = append(nextNodes, node)
	}
	node = Node{a - (intMin(a+b, bMax) - b), intMin(a+b, bMax), OpTransferXY}
	if !beenThere(node, visitedNodes) {
		nextNodes = append(nextNodes, node)
	}
	node = Node{0, b, OpDumpX}
	if !beenThere(node, visitedNodes) {
		nextNodes = append(nextNodes, node)
	}
	node = Node{a, 0, OpDumpY}
	if !beenThere(node, visitedNodes) {
		nextNodes = append(nextNodes, node)
	}
	// create a list of next paths
	for i := range nextNodes {
		temp := make([]Node, len(path)+1)
		copy(temp, path)
		temp[len(temp)-1] = nextNodes[i]
		result = append(result, temp)
	}

	return result
}
