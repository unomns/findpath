package algorithms

import (
	"fmt"
	"unomns/findpath/internal/model"
)

type AStarNode struct {
	Y     int
	X     int
	GCost int // Cost from the start node to this node
	HCost int // Heuristic cost from this node to the end node
	FCost int // Total cost (GCost + HCost)

	Parent *AStarNode
}

type Astar struct {
}

func (a *Astar) Name() string {
	return "A* Search Algorithm"
}

func (a *Astar) Find(m model.GameMap, p model.Player) []model.Node {
	fmt.Printf("Player #%d finding path.. map lenght: %d, map width: %d\n", p.ID, m.Height, m.Width)

	curY := p.StartY
	curX := p.StartX

	fmt.Printf("Start coords: %d %d\n", curY, curX)

	if m.Grid[curY][curX] > 0 {
		fmt.Printf("Wrong position! Only the '0' value is available to moving threw!\n")
	}

	target := AStarNode{Y: p.EndY, X: p.EndX}

	current := &AStarNode{Y: curY, X: curX}
	current.HCost = current.calculateHeuristic(target)
	current.FCost = current.HCost + current.GCost

	var path []model.Node

	loopLimiter := 0

	for current.Y != target.Y && current.X != target.X {
		fmt.Printf("loop %d\n", loopLimiter)
		path = append(path, model.Node{Y: current.Y, X: current.X})

		neighbours := defineNeighboures(current, &m)
		for _, n := range neighbours {
			n.HCost = current.calculateHeuristic(n)
			n.GCost = current.GCost + 1
			n.Parent = current

			// problem here is that i create nodes from scratch
			//   and those that already have been handled and should be skipped - i just have no access to those 'states'
			if m.Grid[n.Y][n.X] > 0 {
				continue
			}

			fmt.Printf("node allowed - Y:%d X:%d\n", n.Y, n.X)
			current = &n
		}

		loopLimiter++
		if loopLimiter > 5 {
			fmt.Println("loop limiter")
			break
		}
	}

	for i, n := range path {
		fmt.Printf("Path found: [%d] [y:%d, x:%d]\n", i, n.Y, n.X)
	}

	return path
}

func defineNeighboures(node *AStarNode, m *model.GameMap) []AStarNode {
	var res []AStarNode

	curY := node.Y
	curX := node.X

	fmt.Printf("curY:%d curX:%d\n", curY, curX)

	if curX > 0 {
		// left neighbor
		fmt.Println("left")
		res = append(res, AStarNode{Y: curY, X: curX - 1})
	}

	if curX < m.Width {
		fmt.Println("right")
		// right neighbor
		res = append(res, AStarNode{Y: curY, X: curX + 1})
	}

	if curY > 0 {
		fmt.Println("top")
		// top neighbor
		res = append(res, AStarNode{Y: curY - 1, X: curX})
	}

	if curY < m.Height {
		fmt.Println("bottom")
		// bottom neighbor
		res = append(res, AStarNode{Y: curY + 1, X: curX})
	}

	return res
}

func (n *AStarNode) calculateHeuristic(to AStarNode) int {
	return abs(n.Y-to.Y) + abs(n.X-to.X) // Manhattan distance or Euclidean distance
}

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}
