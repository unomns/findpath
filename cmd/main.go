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
	file := flag.String("file", "", "Path to the map JSON")
	algorithm := flag.String("algo", "a", "Path finding algorithm (a*, etc)")

	flag.Parse()

	if *file == "" {
		fmt.Println("Usage: myapp -file=map.json -algo=a")
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
		fmt.Printf("Map has invalid format: %v\n", err)
		return
	}

	var algo algorithms.PathFinder

	if algo, err = factory.NewPathFinder(*algorithm); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Algo choosen: '%s'\n", algo.Name())
	pathFindingService := app.NewPathFindingService(algo)

	var wg sync.WaitGroup

	for i, p := range gameMap.Players {
		p.ID = i + 1

		wg.Add(1)
		go func(p model.Player) {
			defer wg.Done()

			pathFindingService.FindPath(gameMap, p)
		}(p)
	}

	wg.Wait()
	fmt.Println("All paths calculated")
}
