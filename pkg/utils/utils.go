package utils

import "strings"

func Find[T any](slice []T, f func(T) bool) (T, bool) {
	for _, elt := range slice {
		if f(elt) {
			return elt, true
		}
	}

	return *new(T), false
}

func SplitStringIntoChunks(msg string, maxLineLength int) []string {
	msgsByNewline := strings.Split(msg, "\n")
	outMessages := []string{}

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > maxLineLength {
			outMessages = append(outMessages, sb.String())
			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	outMessages = append(outMessages, sb.String())
	return outMessages
}
