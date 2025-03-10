package torrent

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     struct {
		Pieces      []byte `bencode:"pieces"`
		PieceLength int    `bencode:"piece length"`
		Length      int    `bencode:"length"`
		Name        string `bencode:"name"`
	} `bencode:"info"`
}
