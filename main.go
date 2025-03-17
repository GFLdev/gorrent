package main

import (
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bittorrent"
	"github.com/GFLdev/gorrent/pkg/logger"
	"os"
)

const UsageMessage = "Usage:" +
	"\n\tgorrent download <torrent-file>" +
	"\n\tgorrent info <torrent-file>" +
	"\n\tgorrent help"

func main() {
	ctx := Context{
		Logger:    logger.NewLogger(nil),
		LocalPeer: &bittorrent.Peer{},
	}

	if len(os.Args) < 2 {
		fmt.Println(UsageMessage)
		os.Exit(2)
	}
	args := os.Args[1:]

	switch args[0] {
	case "info":
		PrintTorrentInfo(args)
	case "help":
		fmt.Println(UsageMessage)
	case "tracker":
		PrintTrackerStatus(args)
	case "download":
		ctx.DownloadTorrent(args)
	default:
		fmt.Println(UsageMessage)
		os.Exit(2)
	}
}
