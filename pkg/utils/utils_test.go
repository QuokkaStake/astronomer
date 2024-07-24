package utils

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func StringOfRandomLength(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestFilter(t *testing.T) {
	t.Parallel()

	array := []string{"true", "false"}
	filtered := Filter(array, func(s string) bool {
		return s == "true"
	})

	assert.Len(t, filtered, 1, "Array should have 1 entry!")
	assert.Equal(t, "true", filtered[0], "Value mismatch!")
}

func TestSplitStringIntoChunksLessThanOneChunk(t *testing.T) {
	t.Parallel()

	str := StringOfRandomLength(10)
	chunks := SplitStringIntoChunks(str, 20)
	assert.Len(t, chunks, 1, "There should be 1 chunk!")
}

func TestSplitStringIntoChunksExactlyOneChunk(t *testing.T) {
	t.Parallel()

	str := StringOfRandomLength(10)
	chunks := SplitStringIntoChunks(str, 10)

	assert.Len(t, chunks, 1, "There should be 1 chunk!")
}

func TestSplitStringIntoChunksMoreChunks(t *testing.T) {
	t.Parallel()

	str := "aaaa\nbbbb\ncccc\ndddd\neeeee\n"
	chunks := SplitStringIntoChunks(str, 10)
	assert.Len(t, chunks, 3, "There should be 3 chunks!")
}
