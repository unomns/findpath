package app_grpc

import (
	"github.com/unomns/findpath/pkg/findpath"
	findpathv1 "github.com/unomns/findpath/protos/gen/findpath"
)

func ToGRPCPaths(paths []*findpath.Path) []*findpathv1.Path {
	res := make([]*findpathv1.Path, len(paths))

	for i, p := range paths {
		fp := &findpathv1.Path{Found: p.Found, PlayerId: p.PlayerID}
		if p.Found {
			fp.Steps = make([]*findpathv1.Node, len(p.Steps))
			for k, s := range p.Steps {
				fp.Steps[k] = &findpathv1.Node{Y: s.Y, X: s.X}
			}
		}

		res[i] = fp
	}

	return res
}

func FromGRPCPlayers(players []*findpathv1.Player) []*findpath.Player {
	res := make([]*findpath.Player, len(players))

	for i, p := range players {
		res[i] = &findpath.Player{
			Start:  findpath.Node{Y: p.Start.Y, X: p.Start.X},
			Target: findpath.Node{Y: p.Target.Y, X: p.Target.X},
		}
	}

	return res
}
