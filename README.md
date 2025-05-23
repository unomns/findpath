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
package main

import (
    "github.com/unomns/findpath/internal/app"
    "github.com/unomns/findpath/internal/model"
    "github.com/unomns/findpath/internal/algorithms/astar"
)

func main() {
    // TODO: show
}
```

## 🌐 Using as a Microservice

### Run Locally
```bash
go run cmd/server/main.go --map=map.json --algo=astar

```

### Run with Docker
```bash
docker build -t unomns/findpath .
docker run -p 8080:8080 unomns/findpath
```