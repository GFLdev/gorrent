package utils

import "strings"

func LPad(text string, n int, char string) string {
	if len(text) >= n {
		return text
	}

	return strings.Repeat(char, n-len(text)) + text
}

func RPad(text string, n int, char string) string {
	if len(text) >= n {
		return text
	}

	return text + strings.Repeat(" ", n-len(text))
}
