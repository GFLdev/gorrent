package bittorrent

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

type Message struct {
	ID   messageID
	Data []byte
}
