package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"unomns/findpath/internal/algorithms"
	"unomns/findpath/internal/app"
	"unomns/findpath/internal/factory"
	"unomns/findpath/internal/model"
)

func main() {
	file := flag.String("file", "map.json", "Path to the map JSON")
	algorithm := flag.String("algo", "a", "Path finding algorithm (a*, etc)")
	debugMode := *flag.Bool("debug", false, "Use debug mode for extended logs")

	flag.Parse()

	if *file == "" {
		fmt.Println("Usage: myapp --file=map.json --algo=a")
		return
	}

	var gameMap model.GameMap
	data, err := os.ReadFile(*file)
	if err != nil {
		fmt.Printf("Read file error: %v\n", err)
		return
	}

	err = json.Unmarshal(data, &gameMap)
	if err != nil {
		fmt.Printf("Game Map has invalid format: %v\n", err)
		return
	}

	for y, arr := range gameMap.Grid {
		for x := range arr {
			gameMap.Map = append(gameMap.Map, model.Node{Y: y, X: x})
		}
	}

	var algo algorithms.PathFinder

	if algo, err = factory.NewPathFinder(*algorithm, debugMode); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if debugMode {
		fmt.Printf("Algo choosen: '%s'\n", algo.Name())
	}
	pathFindingService := app.NewPathFindingService(algo)

	var wg sync.WaitGroup

	for i, p := range gameMap.Players {
		p.ID = i + 1

		wg.Add(1)
		go func() {
			defer wg.Done()

			path := pathFindingService.FindPath(gameMap, &p)

			if path == nil {
				fmt.Println("Target not detected!")

				return
			}

			fmt.Printf("Player #%d Path found:\n", p.ID)
			for k, n := range path {
				fmt.Printf("[%d][y:%d, x:%d]\n", k, n.Y, n.X)
			}
			fmt.Println()
		}()
	}

	wg.Wait()
	fmt.Println("The end")
}
