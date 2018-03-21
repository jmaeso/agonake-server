package agonake

import (
	"errors"
	"fmt"
	"net"

	pt "github.com/jmaeso/agonake-server/pkg/packet_types"
)

type GameManager struct {
	PlayersStore PlayersStorer
}

func (gm *GameManager) RegisterPlayer(msg []string, sender net.Addr) (*Player, error) {
	if len(msg) != 2 {
		return nil, errors.New(pt.Msg + ": Invalid command. Expected: SIGNUP <nick>\n")
	}

	player, err := gm.PlayersStore.NewPlayer(msg, sender)
	if err != nil {
		return nil, err
	}

	return player, nil
}

// RemovePlayer removes the sender's related player and returns true if he was the last.
func (gm *GameManager) RemovePlayer(sender net.Addr) bool {
	gm.PlayersStore.DeletePlayer(sender)

	if len(gm.PlayersStore.GetAllPlayers()) == 0 {
		return true
	}

	return false
}

func (gm *GameManager) GameStateMessage() string {
	players := gm.PlayersStore.GetAllPlayers()
	message := fmt.Sprintf("%s %d", pt.GameState, len(players))

	for _, p := range players {
		message += fmt.Sprintf(" %s %d %d %d %d", p.Nick, p.Color, p.PosX, p.PosY, p.Points)
	}

	return message + "\n"
}
