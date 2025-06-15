package model

type GameMap struct {
	Width   int32     `json:"width"`
	Height  int32     `json:"height"`
	Grid    [][]int32 `json:"grid"`
	Players []Player  `json:"players"`
	Map     []Node    `json:"map"`
}

type Node struct {
	Y int32 `json:"y"`
	X int32 `json:"x"`
}

type Player struct {
	ID     int
	Start  Node
	Target Node
}
