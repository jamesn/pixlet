package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModInt(t *testing.T) {
	// Test positive modulo
	assert.Equal(t, 2, ModInt(5, 3))
	assert.Equal(t, 0, ModInt(6, 3))
	assert.Equal(t, 1, ModInt(7, 3))

	// Test negative modulo (should wrap around)
	assert.Equal(t, 2, ModInt(-1, 3))
	assert.Equal(t, 1, ModInt(-2, 3))
	assert.Equal(t, 0, ModInt(-3, 3))

	// Test with larger numbers
	assert.Equal(t, 3, ModInt(13, 10))
	assert.Equal(t, 7, ModInt(-3, 10))
}

func TestMaxFrameCount(t *testing.T) {
	// Test with empty slice - MaxFrameCount returns 1 even for empty slices
	assert.Equal(t, 1, MaxFrameCount([]Widget{}))

	// Test with single widget
	box := Box{Width: 10, Height: 10}
	assert.Equal(t, 1, MaxFrameCount([]Widget{box}))

	// Test with multiple widgets
	text1 := &Text{Content: "Hello"}
	text2 := &Text{Content: "World"}
	widgets := []Widget{box, text1, text2}

	// All these widgets have 1 frame, so max should be 1
	assert.Equal(t, 1, MaxFrameCount(widgets))
}

func TestMaxFrameCountWithAnimation(t *testing.T) {
	// Create a sequence with multiple children
	// Each child will be shown for 1 frame
	seq := Sequence{
		Children: []Widget{
			Box{Width: 10, Height: 10},
			Box{Width: 20, Height: 20},
			Box{Width: 30, Height: 30},
		},
	}

	// Sequence should have 3 frames (one for each child)
	assert.Equal(t, 3, seq.FrameCount())
	assert.Equal(t, 3, MaxFrameCount([]Widget{seq}))
}
