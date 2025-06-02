package app

import (
	"unomns/findpath/internal/algorithms"
	"unomns/findpath/internal/model"
)

type pathFindingService struct {
	algo algorithms.PathFinder
}

func NewPathFindingService(algo algorithms.PathFinder) *pathFindingService {
	return &pathFindingService{algo: algo}
}

func (s *pathFindingService) FindPath(m model.GameMap, p model.Player) []model.Node {
	return s.algo.Find(m, p)
}
