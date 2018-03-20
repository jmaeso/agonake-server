package agonake

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	agones "agones.dev/agones/sdks/go"
	pt "github.com/jmaeso/agonake-server/pkg/packet_types"
)

type Server struct {
	agonesSDK    *agones.SDK
	conn         net.PacketConn
	healthActive bool
	stop         chan bool
}

// NewServer initializes the following:
//
// udp connection on provided port.
// agones SDK.
// healthckecking system of agones.
// mark agones as ready for incomming connections.
//
// Returns an error if something could not be initialized.
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
		agonesSDK:    agonesSDK,
		conn:         conn,
		healthActive: true,
		stop:         make(chan bool),
	}

	go server.checkingConnectivity()

	err = server.agonesSDK.Ready()
	if err != nil {
		return nil, fmt.Errorf("Could not send ready message: %s", err)
	}

	log.Println("Server started")
	return server, nil
}

func (s *Server) checkingConnectivity() {
	tick := time.Tick(2 * time.Second)
	s.healthActive = true
	for {
		err := s.agonesSDK.Health()
		if err != nil {
			log.Fatalf("Could not send health ping, %v", err)
		}
		select {
		case <-s.stop:
			log.Print("Stopped health pings")
			s.healthActive = false
			return
		case <-tick:
		}
	}
}

// ReceiveAndProcessMsgs is intended to be ran inside the game loop.
// Split the message returned by spaces before passing them for processing.
//
// Returns true if the server has to shutdown.
func (s *Server) ReceiveAndProcessMsgs(b []byte) bool {
	n, sender, err := s.conn.ReadFrom(b)
	if err != nil {
		log.Fatalf("Could not read from udp stream: %s", err)
	}

	txt := strings.Fields(strings.TrimSpace(string(b[:n])))
	log.Printf("Received packet from %s: %s", sender.String(), txt)

	return s.processMsg(txt, sender)
}

func (s *Server) processMsg(msg []string, sender net.Addr) bool {
	switch msg[0] {
	case string(pt.Exit):
		log.Println("Exiting")
		return true

	case string(pt.Unhealthy):
		if s.healthActive {
			close(s.stop)
		}

	default: //echo
		ack := "ACK: " + string(strings.Join(msg[:], " ")) + "\n"
		var err error
		if _, err = s.conn.WriteTo([]byte(ack), sender); err != nil {
			log.Fatalf("Could not write to udp stream: %v", err)
		}
	}

	return false
}

// Stop closes the udp connection and the agones SDK.
func (s *Server) Stop() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("Could not close connection. Err: %s", err)
	}

	err := s.agonesSDK.Shutdown()
	if err != nil {
		return fmt.Errorf("Could not shutdown Agones. Err: %s", err)
	}

	return nil
}
