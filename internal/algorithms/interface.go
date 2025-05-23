package algorithms

import "unomns/findpath/internal/model"

type Position struct {
	X, Y int
}

type PathFinder interface {
	Find(m model.GameMap, p model.Player) model.Path
}
