package algorithms

import (
	"unomns/findpath/internal/model"
)

type Bfs struct{}

func (b *Bfs) Name() string {
	return "Breadth-First Search"
}

func (b *Bfs) Find(m model.GameMap, p *model.Player) []*model.Node {
	// TODO: second
	return make([]*model.Node, 0)
}
