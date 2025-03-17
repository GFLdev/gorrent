package bittorrent

import (
	"encoding/hex"
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bencode"
	"github.com/GFLdev/gorrent/pkg/utils"
)

type TorrentFile struct {
	Announce     string   `bencode:"announce"`
	CreationDate int      `bencode:"creation date"`
	Comment      string   `bencode:"comment"`
	CreatedBy    string   `bencode:"created by"`
	URLList      []string `bencode:"url-list"`
	Info         info     `bencode:"info"`
}

type info struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type TorrentInfo struct {
	Name        string
	TrackerURL  string
	Length      int
	InfoHash    string
	PieceLength int
	PieceHashes []string
}

func TorrentFromFile(torrentPath string) (*TorrentFile, error) {
	torrentFile, err := utils.ReadFile(torrentPath)
	if err != nil {
		return nil, fmt.Errorf("could not read torrent file: %s\n", err.Error())
	}
	torrent := &TorrentFile{}

	err = bencode.Unmarshal(torrentFile, torrent)
	if err != nil {
		return nil, fmt.Errorf("could not parse torrent file: %s\n", err.Error())
	}
	return torrent, nil
}

func (t *TorrentFile) CalculateInfoHash() ([]byte, error) {
	// Bencode info dictionary
	infoBCode, err := bencode.Marshal(&t.Info)
	if err != nil {
		return nil, fmt.Errorf("could not calculate info hash: %w", err)
	}

	// Calculate SHA-1 hash
	return utils.SHA1Encode(infoBCode), nil
}

func GetTorrentInfo(torrentPath string) (TorrentInfo, error) {
	// Read file
	torrentFile, err := utils.ReadFile(torrentPath)
	if err != nil {
		return TorrentInfo{}, fmt.Errorf("could not get torrent info: %w", err)
	}

	// Unmarshal torrent to struct
	var torr TorrentFile
	err = bencode.Unmarshal(torrentFile, &torr)
	if err != nil {
		return TorrentInfo{}, fmt.Errorf("could not get torrent info: %w", err)
	}

	// Check valid pieces
	if len(torr.Info.Pieces)%20 != 0 {
		return TorrentInfo{}, fmt.Errorf("could not get torrent info: invalid pieces")
	}

	idx := 0
	pieceHashes := make([]string, len(torr.Info.Pieces)/20)
	for i := 0; i < len(torr.Info.Pieces); i += 20 {
		pieceHashes[idx] = torr.Info.Pieces[i : i+20]
		idx++
	}

	// Calculate info hash
	infoHash, err := torr.CalculateInfoHash()
	if err != nil {
		return TorrentInfo{}, fmt.Errorf("could not get torrent info: %w", err)
	}

	// Struct info
	torrentInfo := TorrentInfo{
		Name:        torr.Info.Name,
		TrackerURL:  torr.Announce,
		Length:      torr.Info.Length,
		InfoHash:    hex.EncodeToString(infoHash),
		PieceLength: torr.Info.PieceLength,
		PieceHashes: pieceHashes,
	}
	return torrentInfo, nil
}
