package findpath

type Node struct {
	Y int32 `json:"y"`
	X int32 `json:"x"`
}

type Player struct {
	Start  Node `json:"start"`
	Target Node `json:"target"`
}

type Path struct {
	PlayerID string  `json:"player_id"`
	Found    bool    `json:"found"`
	Steps    []*Node `json:"steps"`
}
