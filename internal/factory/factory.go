package factory

import (
	"fmt"
	"unomns/findpath/internal/algorithms"
)

func NewPathFinder(algo string) (algorithms.PathFinder, error) {
	switch algo {
	case "astar":
		return &algorithms.Astar{}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}
