package main

import (
	"flag"
	"fmt"

	"github.com/unomns/findpath/pkg/findpath"
)

func main() {
	file := flag.String("file", "map.example.json", "Path to the map JSON")
	algorithm := flag.String("algo", "a", "Path finding algorithm (a*, etc)")
	debugMode := flag.Bool("debug", false, "Use debug mode for extended logs")

	flag.Parse()

	if *file == "" {
		fmt.Println("Usage: myapp --file=map.json --algo=a")
		return
	}

	service, err := findpath.New(*algorithm, *debugMode)
	if err != nil {
		fmt.Printf("Error! %v\n", err)
		return
	}

	paths, err := service.GetPathFromFile(*file)

	if err != nil {
		fmt.Printf("Error! %v\n", err)
		return
	}

	for i, path := range paths {
		fmt.Printf("Result #%d: %v\n\n", i, path)
	}
}
