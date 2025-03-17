package main

import (
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bittorrent"
	"github.com/GFLdev/gorrent/pkg/logger"
	"net"
	"os"
)

type Context struct {
	Logger    *logger.Logger
	LocalPeer *bittorrent.Peer
}

func PrintTorrentInfo(args []string) {
	if len(args) < 2 {
		fmt.Println(UsageMessage)
		os.Exit(1)
	}
	torrentInfo, err := bittorrent.GetTorrentInfo(args[1])
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	fmt.Println("Name:", torrentInfo.Name)
	fmt.Println("Tracker URL:", torrentInfo.TrackerURL)
	fmt.Println("Length:", torrentInfo.Length)
	fmt.Println("Info Hash:", torrentInfo.InfoHash)
	fmt.Println("Piece Length:", torrentInfo.PieceLength)
	fmt.Println("Piece Hashes:")
	for n, hash := range torrentInfo.PieceHashes {
		if n == 4 {
			fmt.Printf("... (%d more)", len(torrentInfo.PieceHashes)-n)
			break
		}
		fmt.Printf("%x\n", hash)
	}
}

func PrintTrackerStatus(args []string) {
	torrent, err := bittorrent.TorrentFromFile(args[1])
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// Create new local peer
	localPeer := bittorrent.Peer{
		IP:   net.IP{127, 0, 0, 1},
		Port: 6881,
	}

	// Tracker request
	res, err := torrent.AnnounceToTracker(&localPeer)
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

func (ctx *Context) DownloadTorrent(args []string) {
	// Get torrent
	torrent, err := bittorrent.TorrentFromFile(args[1])
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// Create new local peer
	localPeer := bittorrent.Peer{
		IP:   net.IP{127, 0, 0, 1},
		Port: 6881,
	}
	peerId, err := bittorrent.GeneratePeerID()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Tracker request
	res, err := torrent.AnnounceToTracker(&localPeer)
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
		Protocol: bittorrent.BitTorrentProtocol,
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

		fmt.Println("Sending handshake...:", handshake.SerializeHandshake())
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
		fmt.Println("Handshake received!", buf)

		resHandshake, err := bittorrent.DeserializeHandshake(buf)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(resHandshake)
		break
	}
}
