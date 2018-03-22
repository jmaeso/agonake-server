package mock

import (
	"net"

	"github.com/jmaeso/agonake-server/agonake"
)

type PlayersFactory struct{}

func (pf PlayersFactory) Default() *agonake.Player {
	return &agonake.Player{
		Address: &net.IPAddr{
			IP: net.ParseIP("localhost"), Zone: "",
		},
		Nick:   "default",
		Color:  5,
		PosX:   5,
		PosY:   5,
		Points: 100,
	}
}

func (pf PlayersFactory) List(num int) []*agonake.Player {
	players := make([]*agonake.Player, 0)

	for i := 0; i < num; i++ {
		players = append(players, pf.Default())
	}

	return players
}
