package model

type GameMap struct {
	Grid    [][]int  `json:"grid"`
	Players []Player `json:"players"`
}

type Player struct {
	ID     int
	StartX int
	StartY int
	EndX   int
	EndY   int
}

type Path []struct {
	X int `json:"x"`
	Y int `json:"y"`
}
