package algorithms

import (
	"fmt"
	"unomns/findpath/internal/model"
)

type Astar struct {
}

func (a *Astar) Name() string {
	return "A* Search Algorithm"
}

func (a *Astar) Find(m model.GameMap, p model.Player) model.Path {
	fmt.Printf("Player #%d finding path.. map lenght: %d, map width: %d\n", p.ID, m.Height, m.Width)

	curY := p.StartY
	curX := p.StartX

	curPosition := model.Node{Y: curY, X: curX}
	fmt.Printf("Start coords: %d %d\n", curPosition.Y, curPosition.X)

	if m.Grid[curPosition.Y][curPosition.X] > 0 {
		fmt.Printf("Wrong position! Only the '0' value is available to moving threw!\n")
	}

	if curX > 0 {
		leftNeighbor := model.Node{Y: curPosition.Y, X: curPosition.X - 1}
		fmt.Printf("Left neighbor coords: %d %d; value: %d\n",
			leftNeighbor.Y, leftNeighbor.X, m.Grid[leftNeighbor.Y][leftNeighbor.X])
	}

	if curX < m.Width {
		rightNeighbor := model.Node{Y: curPosition.Y, X: curPosition.X + 1}
		fmt.Printf("Right neighbor coords: %d %d; value: %d\n",
			rightNeighbor.Y, rightNeighbor.X, m.Grid[rightNeighbor.Y][rightNeighbor.X])
	}

	if curY > 0 {
		topNeighbor := model.Node{Y: curPosition.Y - 1, X: curPosition.X}
		fmt.Printf("Top neighbor coords: %d %d; value: %d\n",
			topNeighbor.Y, topNeighbor.X, m.Grid[topNeighbor.Y][topNeighbor.X])
	}

	if curY < m.Height {
		buttomtNeighbor := model.Node{Y: curPosition.Y + 1, X: curPosition.X}
		fmt.Printf("Buttom neighbor coords: %d %d; value: %d\n",
			buttomtNeighbor.Y, buttomtNeighbor.X, m.Grid[buttomtNeighbor.Y][buttomtNeighbor.X])
	}

	return model.Path{}
}
