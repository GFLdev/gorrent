package bittorrent

// messageID represents a unique identifier for a specific type of message in a communication protocol.
type messageID uint8

const (
	Choke messageID = iota
	Unchoke
	Interested
	NotInterested
	Have
	Bitfield
	Request
	Piece
	Cancel
)

// Message represents a basic structure for communication containing an identifier and associated data.
type Message struct {
	ID   messageID
	Data []byte
}
