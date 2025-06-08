package algorithms

import (
	"context"
	"fmt"
	"slices"
	"sort"
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
	debugMode bool
	l         map[string][]string
}

func NewAstar(d bool) *Astar {
	return &Astar{debugMode: d, l: make(map[string][]string)}
}

func (a *Astar) Name() string {
	return "A* Search Algorithm"
}

var mutex sync.RWMutex

func (a *Astar) debug(n *AStarNode, msg string) {
	if !a.debugMode {
		return
	}

	if n != nil {
		msg = fmt.Sprintf("node [y:%d,x:%d] | %s", n.Y, n.X, msg)
	}

	// fmt.Println(msg)

	var k string
	if n != nil {
		k = generateKey(n.Y, n.X)
	} else {
		k = "default"
	}

	mutex.RLock()
	_, ok := a.l[k]
	mutex.RUnlock()

	mutex.Lock()
	if !ok {
		a.l[k] = make([]string, 5)
	}

	a.l[k] = append(a.l[k], msg)
	mutex.Unlock()
}

func (a *Astar) PrintDebugLogs() {
	keys := make([]string, len(a.l))

	i := 0
	for k := range a.l {
		keys[i] = k
		i++
	}

	sort.SliceStable(keys, func(i, j int) bool { return keys[i] > keys[j] })

	for _, k := range keys {
		for _, v := range slices.Backward(a.l[k]) {
			fmt.Println(v)
		}
	}
}

func (a *Astar) Find(m model.GameMap, p *model.Player) []*model.Node {
	curY := p.StartY
	curX := p.StartX

	if m.Grid[curY][curX] > 0 {
		a.debug(nil, "Wrong position! Only the '0' value is available to moving threw!")

		return nil
	}

	target := &AStarNode{Y: p.EndY, X: p.EndX}
	current := &AStarNode{Y: curY, X: curX}
	current.HCost = current.calculateHeuristic(target)
	current.FCost = current.HCost + current.GCost

	a.debug(nil, fmt.Sprintf("Player #%d finding path.. map lenght: %d, map width: %d\n", p.ID, m.Height, m.Width))
	a.debug(nil, fmt.Sprintf("Start coords: %d %d", curY, curX))
	a.debug(nil, fmt.Sprintf("Target coords: %d %d\n", p.EndY, p.EndX))

	var path []*model.Node
	skipped := make(map[string]*AStarNode)

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan any)

	go func() {
		defer close(ch)
		path = a.loop(ctx, cancel, m, current, target, path, skipped)
	}()

	<-ch

	if a.debugMode {
		a.PrintDebugLogs()
	}

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

			a.debug(n, "====================\nRecursion proccessing START...\n====================================")

			p := a.loop(ctx, cancel, m, n, target, path, skipped)

			a.debug(n, "====================\nRecursion proccessing RESULT...\n====================================")

			if p != nil {
				a.debug(current, fmt.Sprintf("====================\nNew Current coords | y:%d x:%d\n====================================", current.Y, current.X))
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
			a.debug(current, "ctx.Done received!")
			return nil

		default:
			a.debug(current, fmt.Sprintf("Start loop %d", loopLimiter))

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
				a.debug(current, fmt.Sprintf("New Current coords | y:%d x:%d | len(approved) == 1", approved[0].Y, approved[0].X))
				current = approved[0]
			} else {
				a.debug(current, "========== Falling in recursion!\n")
				ch := a.multibranchHandler(ctx, approved, m, current, target, path, skipped)

				select {
				case <-ctx.Done():
					return nil

				case p, ok := <-ch:
					if !ok {
						return nil
					}

					a.debug(current, "### TARGET DEFINED FROM RECURSION ###")
					return p
				}
			}

			loopLimiter++
			if loopLimiter > 15 {
				a.debug(current, "\n#### Canceled by loop limiter!\n")
				cancel()

				break
			}

			if current.Y != target.Y || current.X != target.X {
				a.debug(current, fmt.Sprintf("End of loop %d | continue", loopLimiter))

				continue
			}

			a.debug(current, "\n###### Target detected successfully!!!\n")
			a.debug(current, fmt.Sprintf("End of loop %d | finish loop", loopLimiter))

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
			a.debug(current, "No neigbours found!")
			ch <- false
			return
		}

		approved = a.calculate(current, &neighbours)

		if len(approved) == 0 {
			a.debug(current, "Neibours not approved!")
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

		a.debug(current, fmt.Sprintf("node allowed - Y:%d X:%d", n.Y, n.X))
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
			a.debug(current, fmt.Sprintf("left node: %d %d", curY, curX-1))
			res = append(res, *n)
		}
	}

	// right neighbor
	if curX < (m.Width - 1) {
		if n := defineNeighbour(curY, curX+1, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("right node: %d %d", curY, curX+1))
			res = append(res, *n)
		}
	}

	// top neighbor
	if curY > 0 {
		if n := defineNeighbour(curY-1, curX, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("top node: %d %d", curY-1, curX))
			res = append(res, *n)
		}
	}

	// bottom neighbor
	if curY < (m.Height - 1) {
		if n := defineNeighbour(curY+1, curX, &m.Grid, skipped); n != nil {
			a.debug(current, fmt.Sprintf("bottom node: %d %d", curY+1, curX))
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
