package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmaeso/agonake-server/agonake"
)

func main() {
	port := flag.String("port", "7654", "The port to listen to udp traffic on")
	flag.Parse()

	playersStore := &agonake.PlayersStore{}

	gameManager := &agonake.GameManager{
		PlayersStore: playersStore,
	}

	server, err := agonake.NewServer(*port)
	if err != nil {
		log.Fatalf("Could not start server. Err: %s", err)
	}

	server.SetManager(gameManager)

	var quit = false
	buffer := make([]byte, 1024)

	for !quit {
		quit = server.ReceiveAndProcessMsgs(buffer)
	}

	if err := server.Stop(); err != nil {
		log.Fatalf("Could not stop server. Err: %s", err)
	}

	os.Exit(0)
}
