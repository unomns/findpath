package factory

import (
	"fmt"
	"unomns/findpath/internal/algorithms"
)

func NewPathFinder(algo string) (algorithms.PathFinder, error) {
	switch algo {
	case "a", "a-star":
		return &algorithms.Astar{}, nil
	case "b", "bfs":
		return &algorithms.Bfs{}, nil
	case "d", "dijkstra":
		return &algorithms.Dijkstra{}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}
