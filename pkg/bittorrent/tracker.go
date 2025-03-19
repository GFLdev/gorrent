package bittorrent

import (
	"encoding/binary"
	"fmt"
	"github.com/GFLdev/gorrent/pkg/bencode"
	"github.com/GFLdev/gorrent/pkg/utils"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

// PeersList represents the response from a tracker, containing a parsed peers list and its interval.
type PeersList struct {
	// Interval specifies the wait time in seconds before the next tracker request.
	Interval int
	// Peers contains the list of available peers provided by the tracker.
	Peers []Peer
}

// successResponse represents a successful response from a tracker with bencoded data.
type successResponse struct {
	// Interval specifies the interval in seconds at which the client should contact the tracker for updates.
	Interval int `bencode:"interval"`
	// Peers contains the list of encoded peers in a compact format (6 bytes each: 4 for IP and 2 for port).
	Peers string `bencode:"peers"`
}

// failedResponse represents an error response from a tracker with bencoded data.
type failedResponse struct {
	// FailureReason contains the tracker-provided message detailing why the request failed.
	FailureReason string `bencode:"failure reason"`
}

// TrackerURL constructs a tracker URL with query parameters based on torrent and peer details.
func (t *TorrentFile) TrackerURL(id []byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", fmt.Errorf("could not parse announce url: %w", err)
	}

	// TODO: implement compact parameter
	hash, err := t.InfoHash()
	params := url.Values{
		"info_hash":  []string{string(hash)},
		"peer_id":    []string{string(id)},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"}, // uploaded nothing yet
		"downloaded": []string{"0"}, // downloaded nothing yet
		"compact":    []string{"1"}, // non-compact not implemented
		"left":       []string{strconv.Itoa(t.Info.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

// FetchTracker sends a GET request to the given tracker URL and retrieves the tracker response as a byte slice.
func (t *TorrentFile) FetchTracker(announceURL string) ([]byte, error) {
	res, err := http.Get(announceURL)
	if err != nil {
		return nil, fmt.Errorf("error on tracker request: %w", err)
	}

	// Read and parse response body
	data, err := utils.ExtractResponseData(res)
	if err != nil {
		return nil, fmt.Errorf("error on tracker request: %w", err)
	}
	return data, nil
}

// ParseTrackerResponse parses the tracker's response data, extracts interval and peer information, and handles errors.
func (t *TorrentFile) ParseTrackerResponse(data []byte, timeout int) (PeersList, error) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	// Check failure
	failed := failedResponse{}
	err := bencode.Unmarshal(data, &failed)
	if err != nil {
		return PeersList{}, fmt.Errorf("could not unmarshal tracker response: %w", err)
	}
	if failed.FailureReason != "" {
		return PeersList{}, fmt.Errorf("tracker request failed: %s", failed.FailureReason)
	}

	// Get interval and peer list
	success := successResponse{}
	err = bencode.Unmarshal(data, &success)
	if err != nil {
		return PeersList{}, fmt.Errorf("could not unmarshal tracker response: %w", err)
	}

	// Check if peers list is valid (6 bytes = 4 for IP and 2 for port)
	if len(success.Peers)%6 != 0 {
		return PeersList{}, fmt.Errorf("received malformed peers list")
	}

	// Parse interval and peer list
	idx := 0
	trackerResponse := PeersList{
		Interval: success.Interval,
		Peers:    make([]Peer, len(success.Peers)/6),
	}
	for i := 0; i < len(success.Peers); i += 6 {
		ip := net.IP(success.Peers[i : i+4])
		port := binary.BigEndian.Uint16([]byte(success.Peers[i+4 : i+6]))
		trackerResponse.Peers[idx] = *NewPeer(ip, port, timeout)
		idx++
	}
	return trackerResponse, nil
}
