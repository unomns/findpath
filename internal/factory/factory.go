package factory

import (
	"fmt"
	"github.com/unomns/findpath/internal/algorithms"
)

func NewPathFinder(algo string, debugMode bool) (algorithms.PathFinder, error) {
	switch algo {
	case "a", "a-star":
		return algorithms.NewAstar(debugMode), nil
	case "b", "bfs":
		return &algorithms.Bfs{}, nil
	case "d", "dijkstra":
		return &algorithms.Dijkstra{}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}
