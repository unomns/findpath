# 🧭 FindPath

A high-performance, stateless pathfinding microservice (and Go package) designed for real-time, parallel path calculations in 2D maps. Built for games and simulations that need fast, concurrent routing for multiple players.

---

## 🚀 Features

- 🧠 Pluggable strategy pattern (A*, Dijkstra, etc.)
- 🧵 Parallel pathfinding using goroutines (1 per player)
- 🧩 Use as a Go package or standalone service
- 🌐 Run as an HTTP/gRPC microservice
- 🧪 Easy to test and integrate into other systems

---

## 📦 Use as a Go Package

Install it in any Go project:

```bash
go get github.com/unomns/findpath
```

### Example usage:

```go
import (
    "github.com/unomns/findpath/pkg/findpath"
)

func main() {
    service, _ := findpath.New(findpath.AlgoAStar, true)

    paths, _ := service.GetPathFromFlatGrid(5, 5, []int32{...}, []*findpath.Player{...})

    fmt.Println(paths)
}
```

## 🌐 Using as a Microservice

### Run Locally
```bash
make cli
./bin/findpath-cli --debug
```
