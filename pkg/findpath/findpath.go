package findpath

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/unomns/findpath/internal/algorithms"
	"github.com/unomns/findpath/internal/app"
	"github.com/unomns/findpath/internal/factory"
	"github.com/unomns/findpath/internal/model"
)

type FindPathService struct {
	algo  string
	debug bool
}

type Pathfinder interface {
	GetPathFromFile(jsonFilename string) ([]*Path, error)
	GetPathFromFlatGrid(width int32, height int32, grid []int32, players []*Player) ([]*Path, error)
}

const (
	AlgoAStar = "a-star"
)

func New(algo string, debug bool) (*FindPathService, error) {
	if _, err := factory.NewPathFinder(algo, debug); err != nil {
		return nil, fmt.Errorf("invalid algorithm: %w", err)
	}
	return &FindPathService{algo: algo, debug: debug}, nil
}

func (fps *FindPathService) SetAStarSearchingAlgorithm() {
	fps.algo = AlgoAStar
}

func (fps *FindPathService) GetPathFromFlatGrid(width int32, height int32, grid []int32, players []*Player) ([]*Path, error) {
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

	return fps.computePaths(&gameMap)
}

func (fps *FindPathService) GetPathFromFile(jsonFilename string) ([]*Path, error) {
	data, err := os.ReadFile(jsonFilename)
	if err != nil {
		return nil, fmt.Errorf("read file error: %v", err)
	}

	var gameMap model.GameMap
	err = json.Unmarshal(data, &gameMap)
	if err != nil {
		return nil, fmt.Errorf("file has invalid format: %v", err)
	}

	return fps.computePaths(&gameMap)
}

func (fps *FindPathService) computePaths(gameMap *model.GameMap) ([]*Path, error) {
	paths := make([]*Path, len(gameMap.Players))

	var algo algorithms.PathFinder
	var err error

	if algo, err = factory.NewPathFinder(fps.algo, fps.debug); err != nil {
		return nil, err
	}

	if fps.debug {
		log.Printf("Algo choosen: '%s'\n", algo.Name())
		log.Println("--------Map Grid---------")
		for y := 0; y < int(gameMap.Height); y++ {
			fmt.Printf("[%d]", y)
			for x := 0; x < int(gameMap.Width); x++ {
				fmt.Printf(" %d", gameMap.Grid[y][x])
			}
			fmt.Println()
		}
		log.Println("-------------------------")
	}
	pathFindingService := app.NewPathFindingService(algo)

	var wg sync.WaitGroup

	for i, p := range gameMap.Players {
		p.ID = i + 1
		paths[i] = &Path{PlayerID: strconv.Itoa(i), Found: false}

		wg.Add(1)
		go func() {
			defer wg.Done()

			path := pathFindingService.FindPath(*gameMap, &p)

			if path == nil {
				if fps.debug {
					log.Printf("Player #%d Target not detected!\n", p.ID)
				}

				return
			}

			if fps.debug {
				log.Printf("Player #%d Path found [start:%v][end:%v]:\n", p.ID, p.Start, p.Target)
			}

			paths[i].Found = true
			paths[i].Steps = make([]*Node, len(path))

			for k, n := range path {
				if fps.debug {
					log.Printf("[%d] %v\n", k, *n)
				}
				paths[i].Steps[k] = &Node{Y: n.Y, X: n.X}
			}
			log.Println()
		}()
	}

	wg.Wait()

	return paths, nil
}
