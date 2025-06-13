package algorithms

import (
	"container/heap"
	"fmt"
	"slices"
	"sort"
	"sync"
	"unomns/findpath/internal/model"
)

type AStarNode struct {
	coords model.Node
	gCost  int // Cost from the start node to this node
	hCost  int // Heuristic cost from this node to the end node
	fCost  int // Total cost (GCost + HCost)

	parent *AStarNode
	index  int
}

type PriorityQueue []*AStarNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].fCost < pq[j].fCost
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

var mutex sync.RWMutex

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

	target := &AStarNode{coords: model.Node{Y: p.EndY, X: p.EndX}}

	current := &AStarNode{coords: model.Node{Y: curY, X: curX}}
	current.hCost = current.calculateHeuristic(target)
	current.fCost = current.hCost + current.gCost

	heap.Push(&pq, current)

	finalNode := a.loop(m, target, &pq, skipped)

	if a.debugMode {
		a.printDebugLogs()
	}

	if finalNode == nil {
		return nil
	}

	var path []*model.Node
	for n := finalNode; n != nil; n = n.parent {
		path = append(path, &n.coords)
	}

	slices.Reverse(path)

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
		a.debug(current, fmt.Sprintf("[loop:%d] New Current coords | %v", loopCounter, current.coords))

		k := generateKey(current.coords.Y, current.coords.X)
		skipped[k] = current

		neighbours := current.neigbours(&m, skipped)
		if len(neighbours) == 0 {
			a.debug(current, "No neigbours found!")

			continue
		}

		for _, n := range neighbours {
			n.calculate(current)
			heap.Push(pq, n)

			if n.coords.Y == target.coords.Y && n.coords.X == target.coords.X {
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

	curY := n.coords.Y
	curX := n.coords.X

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
	_, ok := skipped[generateKey(y, x)]
	if ok || (*grid)[y][x] > 0 {
		return nil
	}

	return &AStarNode{coords: model.Node{Y: y, X: x}}
}

func generateKey(y int, x int) string {
	return fmt.Sprintf("%d-%d", y, x)
}

func (n *AStarNode) calculate(parent *AStarNode) {
	n.hCost = parent.calculateHeuristic(n)
	n.gCost = parent.gCost + 1
	n.parent = parent
}

func (n *AStarNode) calculateHeuristic(to *AStarNode) int {
	return abs(n.coords.Y-to.coords.Y) + abs(n.coords.X-to.coords.X) // Manhattan distance or Euclidean distance
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
		msg = fmt.Sprintf("node %v | %s", n.coords, msg)
	}

	var k string
	if n != nil {
		k = generateKey(n.coords.Y, n.coords.X)
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
