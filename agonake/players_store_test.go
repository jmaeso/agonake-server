package agonake_test

import (
	"net"
	"testing"

	"github.com/jmaeso/agonake-server/agonake"

	"github.com/jmaeso/agonake-server/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPlayer(t *testing.T) {
	assert := assert.New(t)
	playersFactory := mock.PlayersFactory{}

	var (
		existingPlayers []*agonake.Player
		playersStore    *agonake.PlayersStore
	)

	setup := func() {
		playersStore = &agonake.PlayersStore{
			Players: existingPlayers,
		}
	}

	t.Run("When new player can be created", func(*testing.T) {
		existingPlayers = playersFactory.List(2)
		setup()

		t.Run("it returns a pointer to the created player", func(*testing.T) {
			addr := &net.IPAddr{IP: net.ParseIP("localhost"), Zone: ""}
			player, err := playersStore.NewPlayer([]string{"SIGNUP", "john"}, addr)

			assert.NoError(err)
			assert.NotNil(player)
			assert.Equal(addr, player.Address)
			assert.Equal("john", player.Nick)
		})
	})

	t.Run("When server is already full", func(*testing.T) {
		existingPlayers = playersFactory.List(agonake.MaxPlayers)
		setup()

		t.Run("it returns a pointer to the created player", func(*testing.T) {
			addr := &net.IPAddr{IP: net.ParseIP("localhost"), Zone: ""}
			player, err := playersStore.NewPlayer([]string{"SIGNUP", "john"}, addr)

			assert.Error(err)
			assert.Nil(player)
			assert.Equal(agonake.ErrGameFull, err)
			assert.Equal(agonake.MaxPlayers, len(playersStore.Players))
		})
	})
}

func TestSetNick(t *testing.T) {
	assert := assert.New(t)
	playersFactory := mock.PlayersFactory{}

	var (
		existingPlayers []*agonake.Player
		playersStore    *agonake.PlayersStore
	)

	setup := func() {
		playersStore = &agonake.PlayersStore{
			Players: existingPlayers,
		}
	}

	t.Run("When the nick is not already taken", func(*testing.T) {
		existingPlayers = playersFactory.List(2)
		setup()

		t.Run("it returns the same nick", func(*testing.T) {
			nick := playersStore.SetNick("unique")

			assert.NotZero(nick)
			assert.Equal("unique", nick)
		})
	})

	t.Run("When nick is already chosen", func(*testing.T) {
		existingPlayers = playersFactory.List(1)
		setup()

		t.Run("it returns same nick and a consecutive integer", func(*testing.T) {
			nick := playersStore.SetNick("default")

			assert.NotZero(nick)
			assert.Equal("default1", nick)
		})
	})

	t.Run("When nick is chosen multiple times", func(*testing.T) {
		existingPlayers = playersFactory.List(3)
		existingPlayers[1].Nick = "default1"
		existingPlayers[2].Nick = "default2"
		setup()

		t.Run("it returns same nick and a consecutive integer", func(*testing.T) {
			nick := playersStore.SetNick("default")

			assert.NotZero(nick)
			assert.Equal("default3", nick)
		})
	})
}

func TestDeletePlayer(t *testing.T) {
	assert := assert.New(t)
	playersFactory := mock.PlayersFactory{}

	var (
		existingPlayers []*agonake.Player
		playersStore    *agonake.PlayersStore
	)

	setup := func() {
		playersStore = &agonake.PlayersStore{
			Players: existingPlayers,
		}
	}

	t.Run("When new player is deleted", func(*testing.T) {
		initialPlayers := 2
		existingPlayers = playersFactory.List(initialPlayers)
		setup()

		t.Run("number of players is decreased", func(*testing.T) {
			addr := &net.IPAddr{IP: net.ParseIP("localhost"), Zone: ""}
			playersStore.DeletePlayer(addr)

			assert.NotEmpty(playersStore.Players)
			assert.NotNil(playersStore.Players)
			assert.Equal(initialPlayers-1, len(playersStore.Players))
		})
	})

	t.Run("When player was the last one", func(*testing.T) {
		existingPlayers = playersFactory.List(1)
		setup()

		t.Run("there are no players left in the store", func(*testing.T) {
			addr := &net.IPAddr{IP: net.ParseIP("localhost"), Zone: ""}
			playersStore.DeletePlayer(addr)

			assert.NotNil(playersStore.Players)
			assert.Empty(playersStore.Players)
		})
	})
}
