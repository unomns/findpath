package model

type GameMap struct {
	Width   int      `json:"width"`
	Height  int      `json:"height"`
	Grid    [][]int  `json:"grid"`
	Players []Player `json:"players"`
	Map     []Node   `json:"map"`
}

type Node struct {
	X int `json:"x"`
	Y int `json:"y"`
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
