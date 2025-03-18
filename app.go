package main

import (
	"flag"
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bittorrent"
	"github.com/GFLdev/gorrent/pkg/logger"
	"net"
	"os"
)

type App struct {
	Logger       *logger.Logger
	LocalPeer    *bittorrent.Peer
	Args         []string
	PeersTimeout int
}

func NewApp(lgr *logger.Logger, localPeer *bittorrent.Peer) *App {
	// Get OS args
	CheckArgsLength(os.Args, 2, 999)
	args := os.Args[1:]

	// Check for timeout flag
	timeout := flag.Int("t", 0, "Timeout in seconds")
	flag.Parse()
	println(*timeout)
	return &App{
		Logger:       lgr,
		LocalPeer:    localPeer,
		Args:         args,
		PeersTimeout: *timeout,
	}
}

func (app *App) TorrentInfo() {
	CheckArgsLength(app.Args, 2, 2)

	torrent, err := bittorrent.TorrentFromFile(app.Args[1])
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	torrentInfo, err := bittorrent.GetTorrentMetadata(torrent)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	fmt.Println(torrentInfo.String())
}

func (app *App) TrackerStatus() {
	torrent, err := bittorrent.TorrentFromFile(app.Args[1])
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// Create new local peer
	localPeer := bittorrent.NewPeer(net.IPv4zero, 6881, 0)

	// Tracker request
	res, err := torrent.AnnounceToTracker(localPeer, app.PeersTimeout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Interval: %d\n", res.Interval)
	fmt.Println("Peers:")
	for _, peer := range res.Peers {
		fmt.Println(peer.String())
	}
}

func (app *App) DownloadTorrent() {
	// Get torrent
	torrent, err := bittorrent.TorrentFromFile(app.Args[1])
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// Create new local peer
	localPeer := bittorrent.NewPeer(net.IPv4zero, 6881, 0)
	peerId, err := bittorrent.GeneratePeerID()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Tracker request
	res, err := torrent.AnnounceToTracker(localPeer, app.PeersTimeout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create handshake
	infoHash, err := torrent.CalculateInfoHash()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	handshake := bittorrent.Handshake{
		Protocol: bittorrent.TorrentProtocol,
		InfoHash: [20]byte(infoHash),
		PeerID:   [20]byte(peerId),
	}

	for _, peer := range res.Peers {
		fmt.Println("Connecting to peer:", peer.String())
		err = peer.Connect()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Connected to peer:", peer.String())

		fmt.Println("Sending handshake...:", string(handshake.SerializeHandshake()))
		err = peer.Write(handshake.SerializeHandshake())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Handshake sent!")

		fmt.Println("Reading handshake...")
		buf := make([]byte, 68)
		err = peer.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Handshake received!", string(buf))

		resHandshake, err := bittorrent.DeserializeHandshake(buf)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Peer ID:", string(resHandshake.PeerID[:]))
	}
}
