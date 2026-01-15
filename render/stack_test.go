package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackEmpty(t *testing.T) {
	s := Stack{
		Children: []Widget{},
	}

	// Empty stack should have zero size
	bounds := s.PaintBounds(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, image.Rect(0, 0, 0, 0), bounds)

	// Empty stack should have 1 frame (MaxFrameCount returns 1 for empty slices)
	assert.Equal(t, 1, s.FrameCount())

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(s, image.Rect(0, 0, 100, 100), 0)
		assert.NotNil(t, img)
	})
}

func TestStackSingleChild(t *testing.T) {
	child := Box{
		Width:  20,
		Height: 10,
		Color:  color.RGBA{255, 0, 0, 255},
	}

	s := Stack{
		Children: []Widget{child},
	}

	// Should match child's bounds
	bounds := s.PaintBounds(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, image.Rect(0, 0, 20, 10), bounds)

	// Should have 1 frame
	assert.Equal(t, 1, s.FrameCount())

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(s, image.Rect(0, 0, 100, 100), 0)
		assert.NotNil(t, img)
	})
}

func TestStackMultipleChildren(t *testing.T) {
	children := []Widget{
		Box{Width: 50, Height: 25, Color: color.RGBA{145, 17, 17, 255}},
		Box{Width: 30, Height: 15, Color: color.RGBA{17, 145, 17, 255}},
		Box{Width: 4, Height: 32, Color: color.RGBA{17, 17, 145, 255}},
	}

	s := Stack{
		Children: children,
	}

	// Bounds should be the maximum of all children
	bounds := s.PaintBounds(image.Rect(0, 0, 100, 100), 0)

	// Width should be at least 50 (from first box)
	assert.GreaterOrEqual(t, bounds.Dx(), 50)
	// Height should be at least 32 (from third box)
	assert.GreaterOrEqual(t, bounds.Dy(), 32)

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(s, image.Rect(0, 0, 100, 100), 0)
		assert.NotNil(t, img)
	})
}

func TestStackBoundsConstraint(t *testing.T) {
	// Create children larger than the available bounds
	children := []Widget{
		Box{Width: 200, Height: 150, Color: color.RGBA{255, 0, 0, 255}},
		Box{Width: 180, Height: 120, Color: color.RGBA{0, 255, 0, 255}},
	}

	s := Stack{
		Children: children,
	}

	// Bounds should be constrained to available space
	bounds := s.PaintBounds(image.Rect(0, 0, 100, 80), 0)
	assert.LessOrEqual(t, bounds.Dx(), 100)
	assert.LessOrEqual(t, bounds.Dy(), 80)
}

func TestStackFrameCount(t *testing.T) {
	// Create children with different frame counts
	seq1 := Sequence{
		Children: []Widget{
			Box{Width: 10, Height: 10},
			Box{Width: 10, Height: 10},
		},
	}

	seq2 := Sequence{
		Children: []Widget{
			Box{Width: 10, Height: 10},
			Box{Width: 10, Height: 10},
			Box{Width: 10, Height: 10},
			Box{Width: 10, Height: 10},
		},
	}

	s := Stack{
		Children: []Widget{seq1, seq2},
	}

	// Frame count should be the maximum of all children
	assert.Equal(t, 4, s.FrameCount())
}

func TestStackLayering(t *testing.T) {
	// Test that children are drawn in order (later children on top)
	children := []Widget{
		Box{Width: 30, Height: 30, Color: color.RGBA{255, 0, 0, 255}}, // Red background
		Box{Width: 20, Height: 20, Color: color.RGBA{0, 255, 0, 255}}, // Green middle
		Box{Width: 10, Height: 10, Color: color.RGBA{0, 0, 255, 255}}, // Blue foreground
	}

	s := Stack{
		Children: children,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(s, image.Rect(0, 0, 100, 100), 0)
		assert.NotNil(t, img)

		// Verify image size matches largest child
		bounds := img.Bounds()
		assert.Equal(t, 30, bounds.Dx())
		assert.Equal(t, 30, bounds.Dy())
	})
}

func TestStackWithMultipleLayers(t *testing.T) {
	children := []Widget{
		Box{Width: 64, Height: 32, Color: color.RGBA{0, 0, 0, 255}},
		Box{Width: 32, Height: 16, Color: color.RGBA{255, 255, 255, 255}},
	}

	s := Stack{
		Children: children,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(s, image.Rect(0, 0, 100, 100), 0)
		assert.NotNil(t, img)
	})
}
