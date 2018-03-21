package agonake

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	MaxPlayers = 4
	NumColors  = 10
)

var ErrGameFull = errors.New("game is already full")

type PlayersStore struct {
	Players []*Player
}

type PlayersStorer interface {
	NewPlayer(msg []string, sender net.Addr) (*Player, error)
	GetAllPlayers() []*Player
	DeletePlayer(addr net.Addr)
}

// NewPlayer register a new player to the match.
func (ps *PlayersStore) NewPlayer(msg []string, sender net.Addr) (*Player, error) {
	if len(ps.Players) >= MaxPlayers {
		return nil, ErrGameFull
	}

	rand.Seed(time.Now().Unix())
	player := &Player{
		Address: sender,
		Nick:    ps.SetNick(msg[1]),
		Color:   rand.Intn(NumColors),
	}

	ps.Players = append(ps.Players, player)

	return player, nil
}

func (ps *PlayersStore) SetNick(proposal string) string {
	var (
		valid  = false
		nick   = proposal
		suffix = 1
	)

	if len(ps.Players) != 0 {
		for !valid {
			for _, p := range ps.Players {
				oldSuffix := suffix
				if p.Nick == nick {
					nick = fmt.Sprintf("%s%d", proposal, suffix)
					log.Printf("So I'm proposing %s", nick)
					suffix++
					break
				}
				if oldSuffix == suffix {
					valid = true
				}
			}
		}
	}

	return nick
}

func (ps *PlayersStore) GetAllPlayers() []*Player {
	return ps.Players
}

func (ps *PlayersStore) DeletePlayer(addr net.Addr) {
	for i, p := range ps.Players {
		if p.Address.String() == addr.String() {
			ps.Players[i] = ps.Players[len(ps.Players)-1]
			ps.Players[len(ps.Players)-1] = nil
			ps.Players = ps.Players[:len(ps.Players)-1]
			break
		}
	}
}
