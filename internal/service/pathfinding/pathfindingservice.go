package pathfindingservice

type pathFindingService struct {
	algorightm string
}

func PathFindingService(algorithm string) *pathFindingService {
	return &pathFindingService{
		algorightm: algorithm,
	}
}

const A_STAR_ALGORITHM = "astar"

func (s *pathFindingService) Find() {
	// TODO
}
