//pathfinding package implements pathfinding algorithms such as Dijkstra and A*
package pathfinding

import (
	"fmt"
)

//Defining possible graph elements
const (
	UNKNOWN int = iota - 1
	LAND
	WALL
	START
	STOP
)

type MapDict map[int] map[int] int

//A point is just a set of x, y coordinates with a PointType attached
type Node struct {
	x, y int //Using int for efficiency
	parent *Node
	H int
}

//Create a new node
func NewNode (x, y int) *Node {
	max_int := int(^uint(0) >> 1)
	node := &Node{
		x: x,
		y: y,
		parent: new(Node),
		H: max_int,
	}
	return node
}

//Return string representation of the node
func (self *Node) String() string {
	return fmt.Sprintf("<Node x:%s y:%s addr:%s>", self.x, self.y, &self)
}

//Start, end nodes and a slice of nodes
type Graph struct {
	start, stop *Node
	nodes []*Node
}

//Return a Graph from a map of coordinates (those that are passible)
func NewGraph (map_data MapDict) *Graph {
	var start, stop *Node
	nodes := make([]*Node, len(map_data) + len(map_data[0]))
	for i, row := range map_data {
		for j, _type := range row {
			if _type == LAND || _type == START || _type == STOP {
				node := NewNode(i, j)
				nodes = append(nodes, node)
				if _type == START {
					start = node
				}
				if _type == STOP {
					stop = node
				}
			}
		}
	}
	g := &Graph{
		nodes: nodes,
		start: start,
		stop: stop,
	}
	return g
}

//Get the nodes near some node
func (self *Graph) adjacentNodes(node *Node) []*Node {
	var result []*Node
	for _, n := range self.nodes {
		switch{
		case node.x + 1 == n.x && node.y == n.y:
		case node.x-1 == n.x && node.y == n.y:
		case node.x == n.x && node.y-1 == n.y:
		case node.x == n.x && node.y+1 == n.y:
			result = append(result, n)
		}
	}
	return result
}

func (self *Graph) retracePath(current_node *Node, path []*Node) []*Node {
	var none []*Node
	path = append(path, current_node)
	if current_node.parent != nil {
		return none
	}
	self.retracePath(current_node.parent, path)
	return path
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func removeNode(nodes []*Node, node *Node) []*Node {
	var result []*Node
	for _, n := range nodes{
		if n != node {
			result = append(result, n)
		}
	}
	return result
}

func hasNode(nodes []*Node, node *Node) bool {
	for _, n := range nodes{
		if n == node {
			return true
		}
	}
	return false
}

//Return the node with the minimum H
func minH(nodes []*Node) *Node {
	var result_node *Node
	minH := int(^uint(0) >> 1)
	for _, node := range nodes{
		if node.H <= minH {
			minH = node.H
			result_node = node
		}
	}
	return result_node
}


//A* search algorithm. See http://en.wikipedia.org/wiki/A*_search_algorithm
func Astar(graph *Graph) []*Node {
	var path, openSet, closedSet []*Node
	openSet = append(openSet, graph.start)
	for len(openSet) != 0 {
		//Get the node with the min H
		current := minH(openSet)
		if current == graph.stop {
			return graph.retracePath(current, make([]*Node, 0))
		}
		openSet = removeNode(openSet, current)
		closedSet = append(closedSet, current)
		for _, tile := range graph.adjacentNodes(current) {
			if !hasNode(closedSet, tile) {
				tile.H = (abs(graph.stop.x - tile.x) + abs(graph.stop.y - tile.y)) * 10
				if !hasNode(openSet, tile) {
					openSet = append(openSet, tile)
				}
				tile.parent = current
			}
		}
	}
	return path
}
