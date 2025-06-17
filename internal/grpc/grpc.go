package app_grpc

import (
	"context"
	"errors"
	"fmt"
	"unomns/findpath/pkg/findpath"
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

	service, err := findpath.New(defaultAlgo, debugMode)
	if err != nil {
		return nil, err
	}

	paths, err := service.GetPathFromFlatGrid(width, height, grid, FromGRPCPlayers(players))
	if err != nil {
		return nil, err
	}

	return &findpathv1.PathResponse{
		Path: ToGRPCPaths(paths),
	}, nil
}
