package bittorrent

import (
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port int
}

func (p *Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(p.Port))
}
