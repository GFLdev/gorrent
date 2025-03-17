package bittorrent

import (
	"fmt"
	"github.com/GFLdev/gorrent/pkg/utils"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
	conn *net.Conn
}

func GeneratePeerID() ([]byte, error) {
	pubKey, err := utils.ReadFile(utils.PublicKeyPath)
	if err == nil { // keys found
		return utils.SHA1Encode(pubKey), nil
	}

	// Generate new keys
	keys, err := utils.GenerateKeys()
	if err != nil {
		return nil, fmt.Errorf("could not create new peer: %w", err)
	}
	return utils.SHA1Encode(keys.PublicKey), nil
}

func (p *Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func (p *Peer) Connect() error {
	conn, err := net.Dial("tcp", p.String())
	if err != nil {
		return fmt.Errorf("could not connect to peer %s: %w", p.String(), err)
	}
	p.conn = &conn
	return nil
}

func (p *Peer) Close() {
	if p.conn != nil {
		return
	}

	err := (*p.conn).Close
	if err != nil {
		// TODO: Handle peer closing connection error
		panic(err)
	}
}

func (p *Peer) Read(buf []byte) error {
	_, err := (*p.conn).Read(buf)
	if err != nil {
		return fmt.Errorf("could not read from peer %s: %w", p.String(), err)
	}
	return nil
}

func (p *Peer) Write(buf []byte) error {
	_, err := (*p.conn).Write(buf)
	if err != nil {
		return fmt.Errorf("could not write to peer %s: %w", p.String(), err)
	}
	return nil
}
