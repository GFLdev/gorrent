package utils

import (
	"strconv"
	"strings"
)

// LPad pads the input string `text` on the left with the specified `char` until it reaches the desired length `n`.
func LPad(text string, n int, char string) string {
	if len(text) >= n {
		return text
	}

	return strings.Repeat(char, n-len(text)) + text
}

// RPad pads the input `text` on the right with the specified `char` until it reaches the desired length `n`.
func RPad(text string, n int, char string) string {
	if len(text) >= n {
		return text
	}

	return text + strings.Repeat(" ", n-len(text))
}

// Base16ToHex converts a string in base-16 to its hexadecimal representation.
func Base16ToHex(b string) string {
	var h string
	for i := 0; i < len(b); i++ {
		h += strconv.FormatInt(int64(b[i]), 16)
	}
	return h
}
