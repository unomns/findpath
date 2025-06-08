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
	DebugMode bool
}

func (a *Astar) Name() string {
	return "A* Search Algorithm"
}

var mutex sync.RWMutex

func (a *Astar) debug(n *AStarNode, msg string) {
	if !a.DebugMode {
		return
	}

	if n != nil {
		msg = fmt.Sprintf("node y:%d x:%d | %s\n", n.Y, n.X, msg)
	}

	fmt.Println(msg)
}

func (a *Astar) Find(m model.GameMap, p *model.Player) []*model.Node {
	a.debug(nil, fmt.Sprintf("Player #%d finding path.. map lenght: %d, map width: %d\n", p.ID, m.Height, m.Width))

	curY := p.StartY
	curX := p.StartX

	a.debug(nil, fmt.Sprintf("Start coords: %d %d\n", curY, curX))

	if m.Grid[curY][curX] > 0 {
		a.debug(nil, "Wrong position! Only the '0' value is available to moving threw!")

		return nil
	}

	target := &AStarNode{Y: p.EndY, X: p.EndX}

	a.debug(nil, fmt.Sprintf("Target coords: %d %d\n\n", p.EndY, p.EndX))

	current := &AStarNode{Y: curY, X: curX}
	current.HCost = current.calculateHeuristic(target)
	current.FCost = current.HCost + current.GCost

	var path []*model.Node
	skipped := make(map[string]*AStarNode)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		path = a.loop(ctx, cancel, m, current, target, path, skipped)
	}()

	wg.Wait()

	return path
}

func (a *Astar) multibranchHandler(
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

			a.debug(current, fmt.Sprintf("Recursion proccessing START...\n=====\n", n.Y, n.X))

			p := a.loop(ctx, cancel, m, n, target, path, skipped)

			a.debug(current, fmt.Sprintf("Recursion proccessing RESULT...\n=====\n", n.Y, n.X))

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

func (a *Astar) loop(
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
			a.debug(current, "ctx.Done received!\n\n")
			return nil

		default:
			a.debug(current, fmt.Sprintf("loop %d\n", loopLimiter))

			path = append(path, &model.Node{Y: current.Y, X: current.X})

			k := generateKey(current.Y, current.X)

			mutex.Lock()
			skipped[k] = current
			mutex.Unlock()

			approved := a.definePossibleOptions(ctx, current, &m, skipped)
			if approved == nil {
				return nil
			}

			if len(approved) == 1 {
				current = approved[0]
			} else {
				a.debug(current, "========== Falling in recursion!\n\n")
				ch := a.multibranchHandler(ctx, approved, m, current, target, path, skipped)

				select {
				case <-ctx.Done():
					return nil

				case p, ok := <-ch:
					if !ok {
						return nil
					}

					a.debug(current, "TARGET DEFINED FROM RECURSION\n")
					return p
				}
			}

			loopLimiter++
			if loopLimiter > 15 {
				a.debug(current, "\nBREAKed by loop limiter!\n\n")
				cancel()

				break
			}

			a.debug(current, fmt.Sprintf("Current coords | y:%d x:%d\n\n", current.Y, current.X))
			if current.Y != target.Y || current.X != target.X {
				continue
			}

			a.debug(current, "\nTarget detected successfully!\n\n")
			cancel()
		}

		return path
	}
}

func (a *Astar) definePossibleOptions(
	ctx context.Context,
	current *AStarNode,
	m *model.GameMap,
	skipped map[string]*AStarNode,
) []*AStarNode {
	var approved []*AStarNode

	ch := make(chan bool)

	go func() {
		defer close(ch)

		// TODO: define -> filter -> calculate
		neighbours := a.defineNeighbourNodes(current, m, skipped)
		if len(neighbours) == 0 {
			ch <- false
			return
		}

		approved = a.calculate(current, &neighbours)

		if len(approved) == 0 {
			a.debug(current, fmt.Sprintf("No approved for node: y:%d x:%d!\n", current.Y, current.X))
			ch <- false
			return
		}

		ch <- true
	}()

	select {
	case v := <-ch:
		if !v {
			return nil
		}

		return approved
	case <-ctx.Done():
		return nil
	}
}

func (a *Astar) calculate(
	current *AStarNode,
	neighbours *[]AStarNode,
) []*AStarNode {
	var approved []*AStarNode

	// wg := &sync.WaitGroup{}

	for _, n := range *neighbours {
		n.HCost = current.calculateHeuristic(&n)
		n.GCost = current.GCost + 1
		n.Parent = current

		a.debug(current, fmt.Sprintf("node allowed - Y:%d X:%d\n", n.Y, n.X))
		approved = append(approved, &n)
		// wg.Add(1)

		// go func() {
		// 	n.HCost = current.calculateHeuristic(&n)
		// 	n.GCost = current.GCost + 1
		// 	n.Parent = current

		// 	fmt.Printf("node allowed - Y:%d X:%d\n", n.Y, n.X)
		// 	approved = append(approved, &n)
		// }()
	}

	// wg.Wait()

	return approved
}

func (a *Astar) defineNeighbourNodes(current *AStarNode, m *model.GameMap, skipped map[string]*AStarNode) []AStarNode {
	var res []AStarNode

	curY := current.Y
	curX := current.X

	// left neighbor
	if curX > 0 {
		if n := defineNeighbour(curY, curX-1, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("left node: %d %d\n", curY, curX-1))
			res = append(res, *n)
		}
	}

	// right neighbor
	if curX < (m.Width - 1) {
		if n := defineNeighbour(curY, curX+1, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("right node: %d %d\n", curY, curX+1))
			res = append(res, *n)
		}
	}

	// top neighbor
	if curY > 0 {
		if n := defineNeighbour(curY-1, curX, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("top node: %d %d\n", curY-1, curX))
			res = append(res, *n)
		}
	}

	// bottom neighbor
	if curY < (m.Height - 1) {
		if n := defineNeighbour(curY+1, curX, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("bottom node: %d %d\n", curY+1, curX))
			res = append(res, *n)
		}
	}

	return res
}

func defineNeighbour(y int, x int, grid *[][]int, skipped map[string]*AStarNode) *AStarNode {
	mutex.RLock()
	_, ok := skipped[generateKey(y, x)]
	mutex.RUnlock()

	if ok || (*grid)[y][x] > 0 {
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
