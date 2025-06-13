package algorithms

import (
	"container/heap"
	"fmt"
	"slices"
	"sort"
	"sync"
	"unomns/findpath/internal/model"
)

type PriorityQueue []*AStarNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].FCost < pq[j].FCost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*AStarNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *AStarNode, parent *AStarNode, fcost int) {
	item.Parent = parent
	item.FCost = fcost
	heap.Fix(pq, item.index)
}

type AStarNode struct {
	Y     int
	X     int
	GCost int // Cost from the start node to this node
	HCost int // Heuristic cost from this node to the end node
	FCost int // Total cost (GCost + HCost)

	Parent *AStarNode
	index  int
}

type Astar struct {
	debugMode bool
	logs      map[string][]string
}

func NewAstar(d bool) *Astar {
	return &Astar{debugMode: d, logs: make(map[string][]string)}
}

func (a *Astar) Name() string {
	return "A* Search Algorithm"
}

var (
	mutex sync.RWMutex
)

func (a *Astar) Find(m model.GameMap, p *model.Player) []*model.Node {
	if m.Grid[p.StartY][p.StartX] > 0 {
		a.debug(nil, "Wrong position! Only the '0' value is available to moving threw!")

		return nil
	}

	curY := p.StartY
	curX := p.StartX

	a.debug(nil, fmt.Sprintf("Player #%d finding path.. map lenght: %d, map width: %d\n", p.ID, m.Height, m.Width))
	a.debug(nil, fmt.Sprintf("Start coords: %d %d", curY, curX))
	a.debug(nil, fmt.Sprintf("Target coords: %d %d\n", p.EndY, p.EndX))

	skipped := make(map[string]*AStarNode)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	target := &AStarNode{Y: p.EndY, X: p.EndX}

	current := &AStarNode{Y: curY, X: curX}
	current.HCost = current.calculateHeuristic(target)
	current.FCost = current.HCost + current.GCost

	heap.Push(&pq, current)

	finalNode := a.loop(m, target, &pq, skipped)

	if a.debugMode {
		a.printDebugLogs()
	}

	if finalNode == nil {
		return nil
	}

	var path []*model.Node
	for n := finalNode; n != nil; n = n.Parent {
		path = append(path, &model.Node{Y: n.Y, X: n.X})
	}

	return path
}

func (a *Astar) loop(
	m model.GameMap,
	target *AStarNode,
	pq *PriorityQueue,
	skipped map[string]*AStarNode,
) *AStarNode {
	loopCounter := 0

	for pq.Len() > 0 {
		loopCounter++
		current := heap.Pop(pq).(*AStarNode)
		a.debug(current, fmt.Sprintf("[loop:%d] New Current coords | y:%d x:%d | ", loopCounter, current.Y, current.X))

		k := generateKey(current.Y, current.X)
		mutex.Lock()
		skipped[k] = current
		mutex.Unlock()

		neighbours := current.neigbours(&m, skipped)
		if len(neighbours) == 0 {
			a.debug(current, "No neigbours found!")

			continue
		}

		for _, n := range neighbours {
			n.calculate(current)
			heap.Push(pq, n)

			if n.Y == target.Y && n.X == target.X {
				a.debug(n, "\n###### Target detected successfully!!!\n")
				return n
			}
		}

		a.debug(current, fmt.Sprintf("[loop:%d] End of loop | continue", loopCounter))
	}

	return nil
}

func (n *AStarNode) neigbours(m *model.GameMap, skipped map[string]*AStarNode) []*AStarNode {
	var res []*AStarNode

	curY := n.Y
	curX := n.X

	// left neighbor
	if curX > 0 {
		if neigbour := defineNode(curY, curX-1, &m.Grid, skipped); neigbour != nil {
			res = append(res, neigbour)
		}
	}

	// right neighbor
	if curX < (m.Width - 1) {
		if neigbour := defineNode(curY, curX+1, &m.Grid, skipped); neigbour != nil {
			res = append(res, neigbour)
		}
	}

	// top neighbor
	if curY > 0 {
		if neigbour := defineNode(curY-1, curX, &m.Grid, skipped); neigbour != nil {
			res = append(res, neigbour)
		}
	}

	// bottom neighbor
	if curY < (m.Height - 1) {
		if neigbour := defineNode(curY+1, curX, &m.Grid, skipped); neigbour != nil {
			res = append(res, neigbour)
		}
	}

	return res
}

func defineNode(y int, x int, grid *[][]int, skipped map[string]*AStarNode) *AStarNode {
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

func (n *AStarNode) calculate(parent *AStarNode) {
	n.HCost = parent.calculateHeuristic(n)
	n.GCost = parent.GCost + 1
	n.Parent = parent
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

// DEBUGGING && LOGGING
func (a *Astar) debug(n *AStarNode, msg string) {
	if !a.debugMode {
		return
	}

	if n != nil {
		msg = fmt.Sprintf("node [y:%d,x:%d] | %s", n.Y, n.X, msg)
	}

	var k string
	if n != nil {
		k = generateKey(n.Y, n.X)
	} else {
		k = "default"
	}

	mutex.RLock()
	_, ok := a.logs[k]
	mutex.RUnlock()

	mutex.Lock()
	if !ok {
		a.logs[k] = make([]string, 5)
	}

	a.logs[k] = append(a.logs[k], msg)
	mutex.Unlock()
}

func (a *Astar) printDebugLogs() {
	keys := make([]string, len(a.logs))

	i := 0
	for k := range a.logs {
		keys[i] = k
		i++
	}

	sort.SliceStable(keys, func(i, j int) bool { return keys[i] > keys[j] })

	for _, k := range keys {
		for _, v := range slices.Backward(a.logs[k]) {
			fmt.Println(v)
		}
	}
}
