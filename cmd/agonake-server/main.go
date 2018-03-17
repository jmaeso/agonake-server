package main

import (
	"flag"
	"log"

	"github.com/jmaeso/agonake-server/agonake"
)

// main starts a UDP server that received 1024 byte sized packets at at time
// converts the bytes to a string, and logs the output
func main() {
	port := flag.String("port", "7654", "The port to listen to udp traffic on")
	flag.Parse()

	log.Printf("Creating server")
	server, err := agonake.NewServer(*port)
	if err != nil {
		log.Printf("Could not start server. Err: %s", err)
	}

	server.Loop()
}
