package agonake

import "net"

// Player is the model for representing a user in the game.
type Player struct {
	Address net.Addr
	Nick    string
	Color   int
	PosX    int
	PosY    int
	Points  int
}
