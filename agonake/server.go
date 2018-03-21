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
	gameManager  *GameManager
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

func (s *Server) SetManager(gm *GameManager) {
	s.gameManager = gm
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
	if len(txt) < 1 {
		return false
	}

	return s.processMsg(txt, sender)
}

func (s *Server) processMsg(msg []string, sender net.Addr) bool {
	log.Printf("Processing message from %s: %s", sender.String(), msg)
	var (
		response string
		exit     bool
	)

	switch msg[0] {
	case string(pt.Exit):
		log.Println("Exiting")
		response = string(pt.Msg) + ": bye\n"
		exit = true

	case string(pt.Unhealthy):
		if s.healthActive {
			close(s.stop)
		}

	case string(pt.Signup):
		player, err := s.gameManager.RegisterPlayer(msg, sender)
		if err != nil {
			if err == ErrGameFull {
				if err = s.send(pt.Full+"\n", sender); err != nil {
					log.Fatal(err)
				}
				break
			}
			log.Fatalf("Could not register user. Err %s", err)
			break
		}

		if err = s.send(pt.Hello+" "+player.Nick+"\n", sender); err != nil {
			log.Fatal(err)
		}

		if err = s.broadcast(s.gameManager.GameStateMessage()); err != nil {
			log.Fatal(err)
		}

	case string(pt.Disconnect):
		lastOne := s.gameManager.RemovePlayer(sender)

		if err := s.send(pt.Bye+"\n", sender); err != nil {
			log.Fatal(err)
		}

		if lastOne == true {
			log.Println("Las player left. Shutting down.")
			return true
		}

		if err := s.broadcast(s.gameManager.GameStateMessage()); err != nil {
			log.Fatal(err)
		}

	default: //echo
		response = "ACK: " + string(strings.Join(msg[:], " ")) + "\n"
	}

	if len(response) > 0 {
		s.broadcast(response)
	}

	return exit
}

func (s *Server) send(message string, receiver net.Addr) error {
	var err error
	if _, err = s.conn.WriteTo([]byte(message), receiver); err != nil {
		return fmt.Errorf("Could not send message: %s to player: %s. Err: %s", message, receiver.String(), err)
	}

	return nil
}

func (s *Server) broadcast(message string) error {
	var err error
	for _, p := range s.gameManager.PlayersStore.GetAllPlayers() {
		if err = s.send(message, p.Address); err != nil {
			return err
		}
	}

	return nil
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
