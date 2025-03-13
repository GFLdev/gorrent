package main

import (
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bencode"
	"github.com/GFLdev/gorrent/pkg/bittorrent"
	"github.com/GFLdev/gorrent/pkg/utils"
)

func main() {
	payload, err := utils.ReadFile("debian.torrent")
	if err != nil {
		panic(err)
	}

	torr := new(bittorrent.Torrent)
	err = bencode.Unmarshal(payload, torr)
	if err != nil {
		panic(err)
	}

	_, err = bencode.Marshal(torr)
	if err != nil {
		panic(err)
	}

	fmt.Println(torr.CalculateInfoHash())
}
