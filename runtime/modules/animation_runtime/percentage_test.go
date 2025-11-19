package animation_runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestPercentageFromStarlark_ValidRange(t *testing.T) {
	// Test 0.0
	pct, err := PercentageFromStarlark(starlark.Float(0.0))
	assert.NoError(t, err)
	assert.Equal(t, 0.0, pct.Value)

	// Test 0.5
	pct, err = PercentageFromStarlark(starlark.Float(0.5))
	assert.NoError(t, err)
	assert.Equal(t, 0.5, pct.Value)

	// Test 1.0
	pct, err = PercentageFromStarlark(starlark.Float(1.0))
	assert.NoError(t, err)
	assert.Equal(t, 1.0, pct.Value)
}

func TestPercentageFromStarlark_IntegerValues(t *testing.T) {
	// Test with integers that should convert to floats
	pct, err := PercentageFromStarlark(starlark.MakeInt(0))
	assert.NoError(t, err)
	assert.Equal(t, 0.0, pct.Value)

	pct, err = PercentageFromStarlark(starlark.MakeInt(1))
	assert.NoError(t, err)
	assert.Equal(t, 1.0, pct.Value)
}

func TestPercentageFromStarlark_OutOfRange(t *testing.T) {
	// Test below 0.0
	_, err := PercentageFromStarlark(starlark.Float(-0.1))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid range for percentage")

	// Test above 1.0
	_, err = PercentageFromStarlark(starlark.Float(1.1))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid range for percentage")

	// Test far out of range
	_, err = PercentageFromStarlark(starlark.Float(100.0))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid range for percentage")
}

func TestPercentageFromStarlark_InvalidType(t *testing.T) {
	// Test with string
	_, err := PercentageFromStarlark(starlark.String("0.5"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid type for percentage")

	// Test with None
	_, err = PercentageFromStarlark(starlark.None)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid type for percentage")

	// Test with bool
	_, err = PercentageFromStarlark(starlark.Bool(true))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid type for percentage")
}
