package bittorrent

import (
	"fmt"
	"github.com/GFLdev/gorrent/pkg/utils"
	"net"
	"strconv"
	"time"
)

// DefaultTimeout represents the default duration, in seconds, for waiting before a timeout occurs in operations.
const DefaultTimeout = 10

// Peer represents a network peer with its IP, port, and associated connection.
type Peer struct {
	// IP represents the network address of the peer.
	IP net.IP
	// Port specifies the port number on which the peer is listening for incoming connections.
	Port uint16
	// conn represents the active network connection associated with the peer.
	conn *net.Conn
	// timeout specifies the duration for read/write operations or connection attempts before being terminated.
	timeout time.Duration
}

// GeneratePeerID creates a unique peer ID based on the SHA-1 hash of an existing or newly generated public key.
func GeneratePeerID() ([]byte, error) {
	keys, err := utils.LoadKeys()
	if err == nil {
		return utils.SHA1Encode(keys.PublicKey), nil
	}

	// Generate new keys
	keys, err = utils.GenerateKeys()
	if err != nil {
		return nil, fmt.Errorf("could not create new peer: %w", err)
	}

	// Save new keys
	err = utils.SaveKeys(keys)
	if err != nil {
		return nil, fmt.Errorf("could not save new peer: %w", err)
	}

	return utils.SHA1Encode(keys.PublicKey), nil
}

// NewPeer creates and returns a new Peer instance with the specified IP, port and timeout seconds
// (default: DefaultTimeout).
func NewPeer(ip net.IP, port uint16, timeout int) *Peer {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	return &Peer{
		IP:      ip,
		Port:    port,
		conn:    nil,
		timeout: time.Duration(timeout) * time.Second,
	}
}

// String returns a string representation of the Peer as "IP:port".
func (p *Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

// Connect establishes a TCP connection to the peer within the specified timeout duration and assigns it to the Peer.
func (p *Peer) Connect() error {
	// Check if already has a connection established
	if p.conn != nil {
		return nil
	}

	// New dialer and connection
	dialer := net.Dialer{Timeout: p.timeout}
	conn, err := dialer.Dial("tcp", p.String())
	if err != nil {
		return fmt.Errorf("could not connect to peer %s: %w", p.String(), err)
	}
	p.conn = &conn
	return nil
}

// Close closes the Peer connection if it is active, returning an error if the operation fails.
func (p *Peer) Close() error {
	if p.conn != nil {
		return nil
	}

	err := (*p.conn).Close
	if err != nil {
		return fmt.Errorf("could not close peer %s connection", p.String())
	}
	return nil
}

// Read reads data into the provided byte buffer from the peer's connection and returns an error if the operation fails.
func (p *Peer) Read(buf []byte) error {
	// Check has a connection established
	if p.conn == nil {
		return fmt.Errorf("could not read from peer %s: connection not established", p.String())
	}

	// Set timeout and read
	err := (*p.conn).SetReadDeadline(time.Now().Add(p.timeout))
	if err != nil {
		return fmt.Errorf("could not set read deadline for peer %s: %w", p.String(), err)
	}
	_, err = (*p.conn).Read(buf)
	if err != nil {
		return fmt.Errorf("could not read from peer %s: %w", p.String(), err)
	}
	return nil
}

// Write sends the given byte slice to the peer's connection and returns an error if the operation fails.
func (p *Peer) Write(buf []byte) error {
	// Check has a connection established
	if p.conn == nil {
		return fmt.Errorf("could not write to peer %s: connection not established", p.String())
	}

	// Set timeout and write
	err := (*p.conn).SetWriteDeadline(time.Now().Add(p.timeout))
	if err != nil {
		return fmt.Errorf("could not set write deadline for peer %s: %w", p.String(), err)
	}
	_, err = (*p.conn).Write(buf)
	if err != nil {
		return fmt.Errorf("could not write to peer %s: %w", p.String(), err)
	}
	return nil
}
