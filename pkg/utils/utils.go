package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode"
)

func Map[T, V any](slice []T, f func(T) V) []V {
	result := make([]V, len(slice))

	for index, value := range slice {
		result[index] = f(value)
	}

	return result
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

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

func MaybeRemoveQuotes(source string) string {
	if len(source) > 0 && source[0] == '"' {
		source = source[1:]
	}
	if len(source) > 0 && source[len(source)-1] == '"' {
		source = source[:len(source)-1]
	}

	return source
}

func ParseArgsAsMap(source string) (map[string]string, bool) {
	response := map[string]string{}

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(source, f)

	for _, item := range items {
		if !strings.Contains(item, "=") {
			return response, false
		}
		itemSplit := strings.Split(item, "=")
		response[itemSplit[0]] = MaybeRemoveQuotes(itemSplit[1])
	}

	return response, true
}

func FormatDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}

func FormatPercent(percent float64) string {
	return fmt.Sprintf("%.2f%%", percent*100)
}

func FormatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func FormatSince(since time.Time) string {
	duration := time.Since(since)
	if duration < 0 {
		return fmt.Sprintf("in %s", FormatDuration(-duration))
	} else {
		return fmt.Sprintf("%s ago", FormatDuration(duration))
	}
}
