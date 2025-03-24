package bencode

import (
	"flag"
	"strconv"
)

var SimNumbers = flag.Int("sim", 1000, "number of simulations to run")

const (
	Str int = iota
	Int
	List
	Dict
)

type testData[T interface{}] struct {
	data  T
	bCode string
}

func FormatSeed(seed int64) string {
	return "seed: " + strconv.Itoa(int(seed))
}

func FormatInfo(seed int64, bCode string) string {
	return FormatSeed(seed) + "\nbencode: " + bCode
}
