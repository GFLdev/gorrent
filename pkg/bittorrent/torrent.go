package bittorrent

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bencode"
	"net/url"
	"strconv"
)

type Torrent struct {
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

func (t *Torrent) GetTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", fmt.Errorf("could not parse announce url: %w", err)
	}

	params := url.Values{
		"info_hash":  []string{t.Info.Pieces},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Info.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *Torrent) CalculateInfoHash() (string, error) {
	// Bencode info dictionary
	infoBCode, err := bencode.Marshal(&t.Info)
	if err != nil {
		return "", fmt.Errorf("could not calculate info hash: %w", err)
	}

	// Calculate SHA-1 hash
	h := sha1.New()
	h.Write(infoBCode)
	return hex.EncodeToString(h.Sum(nil)), nil
}
