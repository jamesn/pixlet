package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPieChartBasic(t *testing.T) {
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 255, 255, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}},
		Weights:  []float64{180, 135, 45},
		Diameter: 30,
	}

	// Test PaintBounds
	bounds := pc.PaintBounds(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, image.Rect(0, 0, 30, 30), bounds)

	// Test FrameCount
	assert.Equal(t, 1, pc.FrameCount())

	// Test Paint (should not panic)
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
	})
}

func TestPieChartEqualWeights(t *testing.T) {
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}},
		Weights:  []float64{1, 1, 1},
		Diameter: 40,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 40, 40), 0)
		assert.NotNil(t, img)
		
		// Verify image size
		bounds := img.Bounds()
		assert.Equal(t, 40, bounds.Dx())
		assert.Equal(t, 40, bounds.Dy())
	})
}

func TestPieChartSingleSlice(t *testing.T) {
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}},
		Weights:  []float64{100},
		Diameter: 20,
	}

	// Should paint a full circle
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 20, 20), 0)
		assert.NotNil(t, img)
	})
}

func TestPieChartDifferentSizes(t *testing.T) {
	sizes := []int{10, 20, 30, 50, 100}

	for _, diameter := range sizes {
		pc := PieChart{
			Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}},
			Weights:  []float64{60, 40},
			Diameter: diameter,
		}

		bounds := pc.PaintBounds(image.Rect(0, 0, 200, 200), 0)
		assert.Equal(t, image.Rect(0, 0, diameter, diameter), bounds)

		// Should paint without panicking
		assert.NotPanics(t, func() {
			img := PaintWidget(pc, image.Rect(0, 0, diameter, diameter), 0)
			assert.NotNil(t, img)
		})
	}
}

func TestPieChartMoreColorsThanWeights(t *testing.T) {
	// More colors than weights - should cycle through colors
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}},
		Weights:  []float64{50, 50},
		Diameter: 30,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
	})
}

func TestPieChartFewerColorsThanWeights(t *testing.T) {
	// Fewer colors than weights - should cycle through colors
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}},
		Weights:  []float64{25, 25, 25, 25},
		Diameter: 30,
	}

	// Should paint without panicking (colors will cycle)
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
	})
}

func TestPieChartUnevenWeights(t *testing.T) {
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}},
		Weights:  []float64{10, 30, 60},
		Diameter: 40,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 40, 40), 0)
		assert.NotNil(t, img)
	})
}

func TestPieChartSmallWeights(t *testing.T) {
	// Test with very small weights
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}},
		Weights:  []float64{0.1, 0.9},
		Diameter: 30,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
	})
}

func TestPieChartLargeWeights(t *testing.T) {
	// Test with very large weights
	pc := PieChart{
		Colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}},
		Weights:  []float64{1000, 2000},
		Diameter: 30,
	}

	// Should paint without panicking
	assert.NotPanics(t, func() {
		img := PaintWidget(pc, image.Rect(0, 0, 30, 30), 0)
		assert.NotNil(t, img)
	})
}
