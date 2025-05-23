package algorithms

import (
	"fmt"
	"unomns/findpath/internal/model"
)

type Astar struct {
}

func (a *Astar) Find(m model.GameMap, p model.Player) model.Path {
	fmt.Println("A* algorithm used")
	fmt.Printf("Player #%d finding path..\n", p.ID)
	return model.Path{}
}
