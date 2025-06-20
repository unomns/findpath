package algorithms

import (
	"github.com/unomns/findpath/internal/model"
)

type Dijkstra struct{}

func (d *Dijkstra) Name() string {
	return "Dijkstra's Algorithm"
}

func (d *Dijkstra) Find(m model.GameMap, p *model.Player) []*model.Node {
	// TODO: third
	return make([]*model.Node, 0)
}
