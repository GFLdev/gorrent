package bittorrent

import "fmt"

const BitTorrentProtocol = "BitTorrent protocol"

type Handshake struct {
	Protocol string
	InfoHash [20]byte
	PeerID   [20]byte
}

func DeserializeHandshake(buf []byte) (*Handshake, error) {
	if len(buf) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}
	h := &Handshake{}

	protocolLen := int(buf[0]) // protocol identifier length (19)
	curr := protocolLen + 1
	if len(buf) < curr {
		return nil, fmt.Errorf("invalid handshake: protocol length less than buffer length")
	}
	h.Protocol = string(buf[1:curr]) // protocol identifier string

	curr += 8 // 8 reserved bytes
	if len(buf) < curr+40 {
		return nil, fmt.Errorf("invalid handshake: buffer length less than %d", curr+48)
	}
	curr += copy(h.InfoHash[:], buf[curr:curr+20]) // info sha1 hash
	curr += copy(h.PeerID[:], buf[curr:curr+20])   // peer id
	return h, nil
}

func (h *Handshake) SerializeHandshake() []byte {
	i := 1
	buf := make([]byte, len(h.Protocol)+49)

	buf[0] = byte(len(h.Protocol))      // protocol identifier length (19)
	i += copy(buf[i:], h.Protocol)      // protocol identifier string
	i += copy(buf[i:], make([]byte, 8)) // 8 reserved bytes
	i += copy(buf[i:], h.InfoHash[:])   // info sha1 hash
	i += copy(buf[i:], h.PeerID[:])     // peer id
	return buf
}
