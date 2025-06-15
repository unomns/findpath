package app_grpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"unomns/findpath/internal/algorithms"
	"unomns/findpath/internal/app"
	"unomns/findpath/internal/factory"
	"unomns/findpath/internal/model"
	findpathv1 "unomns/findpath/protos/gen/findpath"
)

type Server struct {
	findpathv1.UnimplementedPathFinderServer
}

func NewServer() *Server {
	return &Server{}
}

var (
	defaultAlgo = "a-star"
	debugMode   = false
)

func (s *Server) Path(
	ctx context.Context,
	req *findpathv1.PathRequest,
) (*findpathv1.PathResponse, error) {
	fmt.Println("Processing..")

	width := req.Width
	height := req.Height
	grid := req.Grid
	players := req.Players

	if len(grid) != int(width*height) {
		return nil, errors.New("grid size does not match width × height")
	}

	gameMap := model.GameMap{
		Grid:    make([][]int32, height),
		Players: make([]model.Player, len(players)),
		Width:   width,
		Height:  height,
		Map:     make([]model.Node, len(grid)),
	}

	var y int32
	for y = 0; y < height; y++ {
		gameMap.Grid[y] = grid[y*width : (y+1)*width]
	}

	for i, p := range players {
		gameMap.Players[i] = model.Player{
			Start:  model.Node{Y: p.Start.Y, X: p.Start.X},
			Target: model.Node{Y: p.Target.Y, X: p.Target.X},
		}
	}

	paths := make([]*findpathv1.Path, len(players))

	var algo algorithms.PathFinder
	var err error

	if algo, err = factory.NewPathFinder(defaultAlgo, debugMode); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, errors.New("algorithm not found")
	}

	if debugMode {
		fmt.Printf("Algo choosen: '%s'\n", algo.Name())
	}
	pathFindingService := app.NewPathFindingService(algo)

	var wg sync.WaitGroup

	for i, p := range gameMap.Players {
		p.ID = i + 1
		paths[i] = &findpathv1.Path{PlayerId: strconv.Itoa(i), Found: false}

		wg.Add(1)
		go func() {
			defer wg.Done()

			path := pathFindingService.FindPath(gameMap, &p)

			if path == nil {
				fmt.Printf("Player #%d Target not detected!\n", p.ID)

				return
			}

			fmt.Printf("Player #%d Path found [start:%v][end:%v]:\n", p.ID, p.Start, p.Target)
			paths[i].Steps = make([]*findpathv1.Node, len(path))
			for k, n := range path {
				fmt.Printf("[%d] %v\n", k, *n)
				paths[i].Steps[k] = &findpathv1.Node{Y: n.Y, X: n.X}
			}
			fmt.Println()
		}()
	}

	wg.Wait()

	fmt.Printf("Algo choosen: '%s'\n", algo.Name())

	return &findpathv1.PathResponse{
		Path: paths,
	}, nil
}
