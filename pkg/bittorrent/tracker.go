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

type TrackerAnnounceResponse struct {
	Interval int
	Peers    []Peer
}

type TrackerSuccessResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type TrackerFailedResponse struct {
	FailureReason string `bencode:"failure reason"`
}

func (t *TorrentFile) BuildTrackerAnnounceURL(peerID []byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", fmt.Errorf("could not parse announce url: %w", err)
	}

	infoHash, err := t.CalculateInfoHash()
	params := url.Values{
		"info_hash":  []string{string(infoHash)},
		"peer_id":    []string{string(peerID)},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"}, // uploaded nothing yet
		"downloaded": []string{"0"}, // downloaded nothing yet
		"compact":    []string{"1"}, // non-compact not implemented
		"left":       []string{strconv.Itoa(t.Info.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) FetchTrackerResponse(peer *Peer, announceURL string) ([]byte, error) {
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

func (t *TorrentFile) ParseTrackerResponse(data []byte) (TrackerAnnounceResponse, error) {
	// Check failure
	failed := TrackerFailedResponse{}
	err := bencode.Unmarshal(data, &failed)
	if err != nil {
		return TrackerAnnounceResponse{}, fmt.Errorf("could not unmarshal tracker response: %w", err)
	}
	if failed.FailureReason != "" {
		return TrackerAnnounceResponse{}, fmt.Errorf("tracker request failed: %s", failed.FailureReason)
	}

	// Get interval and peer list
	success := TrackerSuccessResponse{}
	err = bencode.Unmarshal(data, &success)
	if err != nil {
		return TrackerAnnounceResponse{}, fmt.Errorf("could not unmarshal tracker response: %w", err)
	}

	// Check if peers list is valid (6 bytes = 4 for IP and 2 for port)
	if len(success.Peers)%6 != 0 {
		return TrackerAnnounceResponse{}, fmt.Errorf("received malformed peers list")
	}

	// Parse interval and peer list
	idx := 0
	trackerResponse := TrackerAnnounceResponse{
		Interval: success.Interval,
		Peers:    make([]Peer, len(success.Peers)/6),
	}
	for i := 0; i < len(success.Peers); i += 6 {
		trackerResponse.Peers[idx].IP = net.IP(success.Peers[i : i+4])
		trackerResponse.Peers[idx].Port = binary.BigEndian.Uint16([]byte(success.Peers[i+4 : i+6]))
		idx++
	}
	return trackerResponse, nil
}

func (t *TorrentFile) AnnounceToTracker(peer *Peer) (TrackerAnnounceResponse, error) {
	// Get peer ID
	peerID, err := GeneratePeerID()
	if err != nil {
		return TrackerAnnounceResponse{}, err
	}

	// Tracker URL
	trackerURL, err := t.BuildTrackerAnnounceURL(peerID, peer.Port)
	if err != nil {
		return TrackerAnnounceResponse{}, err
	}

	// Get tracker URL and request data
	data, err := t.FetchTrackerResponse(peer, trackerURL)
	if err != nil {
		return TrackerAnnounceResponse{}, fmt.Errorf("could not get tracker response: %w", err)
	}
	return t.ParseTrackerResponse(data)
}
