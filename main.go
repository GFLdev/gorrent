package main

import (
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bittorrent"
	"github.com/GFLdev/gorrent/pkg/logger"
	"os"
)

// UsageMessage is a constant string that provides usage instructions for the gorrent command-line tool.
const UsageMessage = "Usage:" +
	"\n\tgorrent download <torrent-file>" +
	"\n\tgorrent info <torrent-file>" +
	"\n\tgorrent help"

// CheckArgsLength validates the length of the args slice against specified min and max values.
func CheckArgsLength(args []string, min int, max int) {
	if len(args) < min || len(args) > max {
		fmt.Println(UsageMessage)
		os.Exit(2)
	}
}

func main() {
	app := NewApp(logger.NewLogger(nil), &bittorrent.Peer{})

	switch app.Args[0] {
	case "info":
		app.TorrentInfo()
	case "help":
		fmt.Println(UsageMessage)
	case "tracker":
		app.TrackerStatus()
	case "download":
		app.DownloadTorrent()
	default:
		fmt.Println(UsageMessage)
		os.Exit(2)
	}
}
