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
	// TODO: first
	fmt.Printf("Player #%d finding path..\n", p.ID)
	return model.Path{}
}
