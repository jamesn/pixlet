package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFontList(t *testing.T) {
	fonts := GetFontList()

	// Should have some fonts
	assert.NotEmpty(t, fonts)

	// Should contain some known fonts
	assert.Contains(t, fonts, "tb-8")
	assert.Contains(t, fonts, "tom-thumb")
	assert.Contains(t, fonts, "5x8")
	assert.Contains(t, fonts, "6x10")
}

func TestGetFont(t *testing.T) {
	// Test getting a valid font
	font, err := GetFont("tb-8")
	require.NoError(t, err)
	assert.NotNil(t, font)

	// Test getting another valid font
	font, err = GetFont("tom-thumb")
	require.NoError(t, err)
	assert.NotNil(t, font)

	// Test getting an invalid font
	font, err = GetFont("nonexistent-font")
	assert.Error(t, err)
	assert.Nil(t, font)
	assert.Contains(t, err.Error(), "unknown font")
}

func TestGetFontCaseInsensitive(t *testing.T) {
	// Font names should be case-sensitive
	font1, err1 := GetFont("tb-8")
	font2, err2 := GetFont("TB-8")

	require.NoError(t, err1)
	assert.Error(t, err2) // Should fail because font names are case-sensitive
	assert.NotNil(t, font1)
	assert.Nil(t, font2)
}
