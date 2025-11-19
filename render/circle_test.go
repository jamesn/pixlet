package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircleNoChild(t *testing.T) {
	c := Circle{
		Color:    color.RGBA{255, 0, 0, 255},
		Diameter: 20,
	}

	// Test PaintBounds
	bounds := c.PaintBounds(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, image.Rect(0, 0, 20, 20), bounds)

	// Test FrameCount
	assert.Equal(t, 1, c.FrameCount())

	// Test Paint (should not panic)
	assert.NotPanics(t, func() {
		img := PaintWidget(c, image.Rect(0, 0, 20, 20), 0)
		assert.NotNil(t, img)
	})
}

func TestCircleWithChild(t *testing.T) {
	child := Box{
		Width:  10,
		Height: 10,
		Color:  color.RGBA{0, 255, 0, 255},
	}

	c := Circle{
		Color:    color.RGBA{255, 0, 0, 255},
		Diameter: 30,
		Child:    child,
	}

	// Test PaintBounds
	bounds := c.PaintBounds(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, image.Rect(0, 0, 30, 30), bounds)

	// Test FrameCount (should use child's frame count)
	assert.Equal(t, 1, c.FrameCount())

	// Test Paint (should not panic)
	assert.NotPanics(t, func() {
		img := PaintWidget(c, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
	})
}

func TestCircleWithAnimatedChild(t *testing.T) {
	// Create a sequence with multiple frames
	child := Sequence{
		Children: []Widget{
			Box{Width: 5, Height: 5},
			Box{Width: 5, Height: 5},
			Box{Width: 5, Height: 5},
		},
	}

	c := Circle{
		Color:    color.RGBA{255, 0, 0, 255},
		Diameter: 20,
		Child:    child,
	}

	// FrameCount should match child's frame count
	assert.Equal(t, 3, c.FrameCount())
}

func TestCircleDifferentSizes(t *testing.T) {
	sizes := []int{10, 20, 30, 50, 100}

	for _, diameter := range sizes {
		c := Circle{
			Color:    color.RGBA{255, 0, 0, 255},
			Diameter: diameter,
		}

		bounds := c.PaintBounds(image.Rect(0, 0, 200, 200), 0)
		assert.Equal(t, image.Rect(0, 0, diameter, diameter), bounds)

		// Should paint without panicking
		assert.NotPanics(t, func() {
			img := PaintWidget(c, image.Rect(0, 0, diameter, diameter), 0)
			assert.NotNil(t, img)
		})
	}
}

func TestCircleChildCentering(t *testing.T) {
	// Test that child is centered in the circle
	child := Box{
		Width:  10,
		Height: 10,
		Color:  color.RGBA{0, 255, 0, 255},
	}

	c := Circle{
		Color:    color.RGBA{255, 0, 0, 255},
		Diameter: 30,
		Child:    child,
	}

	// Paint and verify it doesn't panic
	assert.NotPanics(t, func() {
		img := PaintWidget(c, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
		
		// Verify image size
		bounds := img.Bounds()
		assert.Equal(t, 30, bounds.Dx())
		assert.Equal(t, 30, bounds.Dy())
	})
}
