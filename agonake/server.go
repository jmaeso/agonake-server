package agonake

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	agones "agones.dev/agones/sdks/go"
)

type Server struct {
	agonesSDK *agones.SDK
	conn      net.PacketConn
	stop      chan bool
}

func NewServer(port string) (*Server, error) {
	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("Could not start udp server: %s", err)
	}

	agonesSDK, err := agones.NewSDK()
	if err != nil {
		return nil, fmt.Errorf("Could not connect to sdk: %s", err)
	}

	server := &Server{
		agonesSDK: agonesSDK,
		conn:      conn,
		stop:      make(chan bool),
	}

	go server.checkingConnectivity()

	err = server.agonesSDK.Ready()
	if err != nil {
		return nil, fmt.Errorf("Could not send ready message: %s", err)
	}

	return server, nil
}

func (s *Server) checkingConnectivity() {
	tick := time.Tick(2 * time.Second)
	for {
		err := s.agonesSDK.Health()
		if err != nil {
			log.Fatalf("Could not send health ping, %v", err)
		}
		select {
		case <-s.stop:
			log.Print("Stopped health pings")
			return
		case <-tick:
		}
	}
}

func (s *Server) Loop() {
	b := make([]byte, 1024)
	for {
		n, sender, err := s.conn.ReadFrom(b)
		if err != nil {
			log.Fatalf("Could not read from udp stream: %v", err)
		}

		txt := strings.TrimSpace(string(b[:n]))
		log.Printf("Received packet from %v: %v", sender.String(), txt)
		switch txt {
		// shuts down the gameserver
		case "EXIT":
			if err := s.close(); err != nil {
				log.Printf("Could not close connection. Err: %s", err)
			}
			// This tells Agones to shutdown this Game Server
			err := s.agonesSDK.Shutdown()
			if err != nil {
				log.Printf("Could not shutdown Agones. Err: %s", err)
			}
			os.Exit(0)
			break

		// turns off the health pings
		case "UNHEALTHY":
			close(s.stop)
			break
		}

		// echo it back
		ack := "ACK: " + txt + "\n"
		if _, err = s.conn.WriteTo([]byte(ack), sender); err != nil {
			log.Fatalf("Could not write to udp stream: %v", err)
		}
	}
}

func (s *Server) close() error {
	return s.conn.Close()
}
