package bittorrent

import (
	"encoding/hex"
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bencode"
	"github.com/GFLdev/gorrent/pkg/utils"
	"strconv"
)

// TorrentFile represents the metadata structure of a .torrent file parsed according to the BitTorrent specification.
type TorrentFile struct {
	// Announce specifies the primary tracker URL for the torrent.
	Announce string `bencode:"announce"`
	// CreationDate represents the timestamp of when the torrent was created.
	CreationDate int `bencode:"creation date"`
	// Comment holds an optional textual description.
	Comment string `bencode:"comment"`
	// CreatedBy specifies the name and version of the application used to create the torrent file.
	CreatedBy string `bencode:"created by"`
	// URLList specifies an optional list of tracker URLs.
	URLList []string `bencode:"url-list"`
	// Info represents the bencoded "info" dictionary containing essential metadata for the torrent.
	Info struct {
		// Pieces is a concatenated string of SHA-1 hashes.
		Pieces string `bencode:"pieces"`
		// PieceLength specifies the size of each piece in bytes.
		PieceLength int `bencode:"piece length"`
		// Length represents the total size of the file in bytes.
		Length int `bencode:"length"`
		// Name specifies the name of the primary file or the directory name.
		Name string `bencode:"name"`
	} `bencode:"info"`
}

// TorrentMetadata represents metadata information parsed from a torrent file.
type TorrentMetadata struct {
	// Name is the name of the torrent file or content.
	Name string
	// TrackerURL is the URL of the tracker used to announce peers.
	TrackerURL string
	// Length is the total size of the torrent content in bytes.
	Length int
	// InfoHash is a SHA-1 hash of the torrent's info dictionary, used to uniquely identify the torrent.
	InfoHash string
	// PieceLength is the size of each piece in bytes, except possibly the last one.
	PieceLength int
	// PieceHashes contains a list of SHA-1 hashes for each piece of the torrent content.
	PieceHashes []string
}

// TorrentFromFile reads a torrent file from the given path, parses its content, and returns a TorrentFile instance.
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

// CalculateInfoHash generates and returns the SHA-1 hash of the bencoded info dictionary of the torrent.
func (t *TorrentFile) CalculateInfoHash() ([]byte, error) {
	// Bencode info dictionary
	infoBCode, err := bencode.Marshal(&t.Info)
	if err != nil {
		return nil, fmt.Errorf("could not calculate info hash: %w", err)
	}

	// Calculate SHA-1 hash
	return utils.SHA1Encode(infoBCode), nil
}

// GetTorrentMetadata extracts and returns metadata about a torrent.
func GetTorrentMetadata(torr *TorrentFile) (TorrentMetadata, error) {
	// Check valid pieces
	if len(torr.Info.Pieces)%20 != 0 {
		return TorrentMetadata{}, fmt.Errorf("could not get torrent info: invalid pieces")
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
		return TorrentMetadata{}, fmt.Errorf("could not get torrent info: %w", err)
	}

	// Struct info
	torrentInfo := TorrentMetadata{
		Name:        torr.Info.Name,
		TrackerURL:  torr.Announce,
		Length:      torr.Info.Length,
		InfoHash:    hex.EncodeToString(infoHash),
		PieceLength: torr.Info.PieceLength,
		PieceHashes: pieceHashes,
	}
	return torrentInfo, nil
}

// String returns a formatted string representation of the TorrentMetadata, including metadata and piece hash details.
func (tInfo *TorrentMetadata) String() string {
	infoStr := "Name: " + tInfo.Name + "\n" +
		"Tracker URL: " + tInfo.TrackerURL + "\n" +
		"Length: " + strconv.Itoa(tInfo.Length) + "\n" +
		"Info Hash: " + tInfo.InfoHash + "\n" +
		"Piece Length: " + strconv.Itoa(tInfo.PieceLength) + "\n" +
		"Piece Hashes:"

	if len(tInfo.PieceHashes) == 0 {
		infoStr += " (empty)"
	} else if len(tInfo.PieceHashes) == 1 {
		infoStr += " " + utils.Base16ToHex(tInfo.PieceHashes[0])
	} else {
		for n, hash := range tInfo.PieceHashes {
			if n == 5 { // print 5 hashes at maximum
				remain := strconv.Itoa(len(tInfo.PieceHashes) - n)
				infoStr += "\n... (" + remain + " more)"
				break
			}
			infoStr += "\n" + utils.Base16ToHex(hash)
		}
	}
	return infoStr
}
