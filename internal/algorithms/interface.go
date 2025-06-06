package algorithms

import "unomns/findpath/internal/model"

type Position struct {
	X, Y int
}

type PathFinder interface {
	Name() string
	Find(m model.GameMap, p model.Player) []*model.Node
}
