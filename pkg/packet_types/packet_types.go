package packet_types

// PacketType is a type defined for having a list with all the available message types the server accepts.
type PacketType string

const (
	// Exit is the PacketType for closing the connection
	Exit PacketType = "EXIT"

	// Unhealthy is the PacketType for stopping health checking of agones SDK.
	Unhealthy = "UNHEALTHY"

	// Signup is the PacketType for a new client to be registered.
	Signup = "SIGNUP"
)
