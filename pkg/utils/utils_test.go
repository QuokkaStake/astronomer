package utils

import (
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/math"

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

func TestMap(t *testing.T) {
	t.Parallel()

	array := []int64{1, 2, 3}

	filtered := Map(array, func(value int64) int64 {
		return value * 2
	})

	assert.Len(t, filtered, 3)
	assert.Equal(t, int64(2), filtered[0])
	assert.Equal(t, int64(4), filtered[1])
	assert.Equal(t, int64(6), filtered[2])
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

func TestFind(t *testing.T) {
	t.Parallel()

	array := []string{"1", "2", "3"}

	value1, found1 := Find(array, func(s string) bool {
		return s == "1"
	})
	assert.Equal(t, "1", value1)
	assert.True(t, found1)

	value2, found2 := Find(array, func(s string) bool {
		return s == "4"
	})
	assert.Equal(t, "", value2)
	assert.False(t, found2)
}

func TestGroupBy(t *testing.T) {
	t.Parallel()

	type TestType struct {
		Keys  []string
		Value string
	}

	array := []TestType{
		{Keys: []string{"key1"}, Value: "value1"},
		{Keys: []string{"key2"}, Value: "value2"},
		{Keys: []string{"key2"}, Value: "value3"},
		{Keys: []string{"key2", "key3"}, Value: "value4"},
	}

	grouped := GroupBy(array, func(value TestType) []string {
		return value.Keys
	})

	assert.Equal(t, map[string][]TestType{
		"key1": {{Keys: []string{"key1"}, Value: "value1"}},
		"key2": {
			{Keys: []string{"key2"}, Value: "value2"},
			{Keys: []string{"key2"}, Value: "value3"},
			{Keys: []string{"key2", "key3"}, Value: "value4"},
		},
		"key3": {{Keys: []string{"key2", "key3"}, Value: "value4"}},
	}, grouped)
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

func TestFormatDuration(t *testing.T) {
	t.Parallel()

	duration := time.Hour*24 + time.Hour*2 + time.Second*4
	formatted := FormatDuration(duration)

	assert.Equal(t, "1 day 2 hours 4 seconds", formatted)
}

func TestFormatLegacyDec(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "123,456,789.123", FormatDec(math.LegacyMustNewDecFromStr("123456789.123456")))
	assert.Equal(t, "1,234,567.123", FormatDec(math.LegacyMustNewDecFromStr("1234567.123456")))
}

func TestFormatPercent(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "10.00%", FormatPercent(0.1))
}

func TestFormatFloat(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "10.00", FormatFloat(10.00))
}

func TestMaybeRemoveQuotes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "string", MaybeRemoveQuotes("string"))
	assert.Equal(t, "string", MaybeRemoveQuotes("\"string\""))
}

func TestParseArgsAsMap(t *testing.T) {
	t.Parallel()

	_, invalid := ParseArgsAsMap("a b c d")
	assert.False(t, invalid)

	args, valid := ParseArgsAsMap("key1=value1 key2=\"value2\" key3=\"value3 value4\"")
	assert.True(t, valid)
	assert.Equal(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3 value4",
	}, args)
}

func TestFormatSince(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "in 1 minute", FormatSince(time.Now().Add(time.Minute+time.Second)))
	assert.Equal(t, "1 minute ago", FormatSince(time.Now().Add(-time.Minute)))
}
