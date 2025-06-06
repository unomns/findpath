package algorithms

import (
	"context"
	"fmt"
	"sync"
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

var mutex sync.RWMutex

func (a *Astar) Find(m model.GameMap, p model.Player) []*model.Node {
	fmt.Printf("Player #%d finding path.. map lenght: %d, map width: %d\n", p.ID, m.Height, m.Width)

	curY := p.StartY
	curX := p.StartX

	fmt.Printf("Start coords: %d %d\n", curY, curX)

	if m.Grid[curY][curX] > 0 {
		fmt.Printf("Wrong position! Only the '0' value is available to moving threw!\n")
	}

	target := &AStarNode{Y: p.EndY, X: p.EndX}

	fmt.Printf("Target coords: %d %d\n\n", p.EndY, p.EndX)

	current := &AStarNode{Y: curY, X: curX}
	current.HCost = current.calculateHeuristic(target)
	current.FCost = current.HCost + current.GCost

	var path []*model.Node
	skipped := make(map[string]*AStarNode)

	// ch := make(chan *AStarNode)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		path = loop(ctx, cancel, m, current, target, path, skipped)
	}()

	wg.Wait()

	if path == nil {
		fmt.Println("Target not detected!")
		return nil
	}

	for i, n := range path {
		fmt.Printf("Path found: [%d] [y:%d, x:%d]\n", i, n.Y, n.X)
	}

	return path
}

func myultibranchHandler(
	parent context.Context,
	approved []*AStarNode,
	m model.GameMap,
	current *AStarNode,
	target *AStarNode,
	path []*model.Node,
	skipped map[string]*AStarNode,
) chan []*model.Node {
	ch := make(chan []*model.Node)

	ctx, cancel := context.WithCancel(parent)
	// defer cancel()

	wg := &sync.WaitGroup{}

	for _, n := range approved {
		wg.Add(1)

		go func() {
			defer wg.Done()

			fmt.Printf("=====node y:%d x:%d | Recursion proccessing START...\n=====\n", n.Y, n.X)

			p := loop(ctx, cancel, m, n, target, path, skipped)

			fmt.Printf("=====node y:%d x:%d | Recursion proccessing RESULT...\n=====\n", n.Y, n.X)

			if p != nil {
				current = n
				ch <- p
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func loop(
	ctx context.Context,
	cancel context.CancelFunc,
	m model.GameMap,
	current *AStarNode,
	target *AStarNode,
	path []*model.Node,
	skipped map[string]*AStarNode,
) []*model.Node {
	loopLimiter := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("ctx.Done received!\n\n")
			return nil

		default:
			fmt.Printf("loop %d\n", loopLimiter)
			path = append(path, &model.Node{Y: current.Y, X: current.X})

			k := generateKey(current.Y, current.X)

			mutex.Lock()
			skipped[k] = current
			mutex.Unlock()

			neighbours := defineNeighboures(current, &m, skipped)
			if len(neighbours) == 0 {
				return nil
			}

			approved := defineApprovedNodes(neighbours, skipped, current)

			l := len(approved)

			if l == 0 {
				fmt.Printf("No approved for node: y:%d x:%d!\n", current.Y, current.X)
				return nil
			}

			if l == 1 {
				current = approved[0]
			} else {
				fmt.Printf("========== Falling in recursion!\n\n")
				ch := myultibranchHandler(ctx, approved, m, current, target, path, skipped)

				select {
				case <-ctx.Done():
					return nil

				case p, ok := <-ch:
					if !ok {
						return nil
					}

					fmt.Printf("TARGET DEFINED FROM RECURSION\n")
					return p
				}
			}

			loopLimiter++
			if loopLimiter > 15 {
				fmt.Printf("\nBREAKed by loop limiter!\n\n")
				cancel()

				break
			}

			fmt.Printf("Current coords | y:%d x:%d\n\n", current.Y, current.X)
			if current.Y != target.Y || current.X != target.X {
				continue
			}

			fmt.Printf("\nTarget detected successfully!\n\n")
			cancel()
		}

		return path
	}
}

func defineApprovedNodes(
	neighbours []AStarNode,
	skipped map[string]*AStarNode,
	current *AStarNode,
) []*AStarNode {
	var approved []*AStarNode

	for _, n := range neighbours {
		var k string = fmt.Sprintf("%d-%d", n.Y, n.X)
		if _, ok := skipped[k]; ok {
			fmt.Printf("Skipping.. y:%d x:%d\n", n.Y, n.X)
			continue
		}

		n.HCost = current.calculateHeuristic(&n)
		n.GCost = current.GCost + 1
		n.Parent = current

		fmt.Printf("node allowed - Y:%d X:%d\n", n.Y, n.X)
		approved = append(approved, &n)
	}

	return approved
}

func defineNeighboures(node *AStarNode, m *model.GameMap, skipped map[string]*AStarNode) []AStarNode {
	var res []AStarNode

	curY := node.Y
	curX := node.X

	fmt.Printf("curY:%d curX:%d\n", curY, curX)

	// left neighbor
	if curX > 0 {
		if n := defineNeighbour(curY, curX-1, &m.Grid, skipped); n != nil {
			fmt.Printf("left node: %d %d\n", curY, curX-1)
			res = append(res, *n)
		}
	}

	// right neighbor
	if curX < (m.Width - 1) {
		if n := defineNeighbour(curY, curX+1, &m.Grid, skipped); n != nil {
			fmt.Printf("right node: %d %d\n", curY, curX+1)
			res = append(res, *n)
		}
	}

	// top neighbor
	if curY > 0 {
		if n := defineNeighbour(curY-1, curX, &m.Grid, skipped); n != nil {
			fmt.Printf("top node: %d %d\n", curY-1, curX)
			res = append(res, *n)
		}
	}

	// bottom neighbor
	if curY < (m.Height - 1) {
		if n := defineNeighbour(curY+1, curX, &m.Grid, skipped); n != nil {
			fmt.Printf("bottom node: %d %d\n", curY+1, curX)
			res = append(res, *n)
		}
	}

	return res
}

func defineNeighbour(y int, x int, grid *[][]int, skipped map[string]*AStarNode) *AStarNode {
	if _, ok := skipped[generateKey(y, x)]; ok || (*grid)[y][x] > 0 {
		return nil
	}

	return &AStarNode{Y: y, X: x}
}

func generateKey(y int, x int) string {
	return fmt.Sprintf("%d-%d", y, x)
}

func (n *AStarNode) calculateHeuristic(to *AStarNode) int {
	return abs(n.Y-to.Y) + abs(n.X-to.X) // Manhattan distance or Euclidean distance
}

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}
