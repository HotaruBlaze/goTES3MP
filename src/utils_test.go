package main

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_utils_appendIfMissing(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	// Append a missing element
	newSlice := appendIfMissing(slice, "durian")
	expectedSlice := []string{"apple", "banana", "cherry", "durian"}
	assert.Equal(t, expectedSlice, newSlice)

	// Don't append an existing element
	newSlice = appendIfMissing(slice, "banana")
	assert.Equal(t, slice, newSlice)
}

func Test_utils_ToHexInt(t *testing.T) {
	n := big.NewInt(123456)

	hex := toHexInt(n)
	expectedHex := "01E240"
	assert.Equal(t, expectedHex, hex)
}

func Test_utils_FindinArray(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	// Find an existing element
	index, found := FindinArray(slice, "banana")
	assert.True(t, found)
	assert.Equal(t, 1, index)

	// Try to find a missing element
	index, found = FindinArray(slice, "durian")
	assert.False(t, found)
	assert.Equal(t, -1, index)
}

func Test_utils_RemoveRGBHex(t *testing.T) {
	message := "The color is #FF0000"
	expectedMessage := "The color is "
	newMessage := removeRGBHex(message)
	assert.Equal(t, expectedMessage, newMessage)

	// Test with multiple RGB Hex values
	message = "The colors are #FF0000, #00FF00, and #0000FF"
	expectedMessage = "The colors are , , and "
	newMessage = removeRGBHex(message)
	assert.Equal(t, expectedMessage, newMessage)
}

func Test_utils_FilterDiscordEmotes(t *testing.T) {
	message := "This is a <:emote_name:123456> emote"
	expectedMessage := "This is a :emote_name: emote"
	newMessage := filterDiscordEmotes(message)
	assert.Equal(t, expectedMessage, newMessage)

	// Test with multiple emotes
	message = "These are <:emote1:123> and <:emote2:456> emotes"
	expectedMessage = "These are :emote1: and :emote2: emotes"
	newMessage = filterDiscordEmotes(message)
	assert.Equal(t, expectedMessage, newMessage)
}
