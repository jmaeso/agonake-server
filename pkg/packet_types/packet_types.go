package packet_types

// PacketType is a type defined for having a list with all the available message types the server accepts.
type PacketType string

const (
	// Exit is the PacketType for closing the connection
	Exit PacketType = "EXIT"

	// Unhealthy is the PacketType for stopping health checking of agones SDK.
	Unhealthy = "UNHEALTHY"

	// Msg is the PacketType for sending text message to clients.
	Msg = "MSG"

	// Signup is the PacketType for a new client to be registered.
	//
	// Expected command: SIGNUP <nick_proposal>
	Signup = "SIGNUP"

	// Full is the PacketType used when requesting signup but server is full.
	Full = "FULL"

	// Hello is the PacketType for telling a client he got registered.
	//
	// Expected response: HELLO <assigned_nick>
	Hello = "HELLO"

	// GameState is the PacketType for updating the game match.
	//
	// Expected response: GS <num_players> $numPlayers_times[<nick> <color> <x_pos> <y_pos> <points>]
	GameState = "GS"

	// Disconnect is the PacketType for requesting disconnection.
	Disconnect = "DISCONNECT"

	// Bye is the PacketType for responding to client a successful disconnection.
	Bye = "BYE"
)
