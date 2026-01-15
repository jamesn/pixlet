package animation_runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render/animation"
)

func TestCurveFromStarlark_EmptyString(t *testing.T) {
	curve, err := CurveFromStarlark(starlark.String(""))
	assert.NoError(t, err)
	_, ok := curve.(animation.LinearCurve)
	assert.True(t, ok)
}

func TestCurveFromStarlark_ValidCurveString(t *testing.T) {
	curve, err := CurveFromStarlark(starlark.String("linear"))
	assert.NoError(t, err)
	_, ok := curve.(animation.LinearCurve)
	assert.True(t, ok)
}

func TestCurveFromStarlark_InvalidCurveString(t *testing.T) {
	curve, err := CurveFromStarlark(starlark.String("invalid_curve"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid curve string")
	_, ok := curve.(animation.LinearCurve)
	assert.True(t, ok)
}

func TestCurveFromStarlark_Function(t *testing.T) {
	code := `
def curve_func(t):
    return t * 2
`
	thread := &starlark.Thread{}
	globals, err := starlark.ExecFile(thread, "test.star", code, nil)
	assert.NoError(t, err)

	fn := globals["curve_func"].(*starlark.Function)
	curve, err := CurveFromStarlark(fn)
	assert.NoError(t, err)

	customCurve, ok := curve.(animation.CustomCurve)
	assert.True(t, ok)
	assert.Equal(t, fn, customCurve.Function)
}

func TestCurveFromStarlark_InvalidFunctionParams(t *testing.T) {
	code := `
def bad_func():
    return 0.5
`
	thread := &starlark.Thread{}
	globals, err := starlark.ExecFile(thread, "test.star", code, nil)
	assert.NoError(t, err)

	fn := globals["bad_func"].(*starlark.Function)
	_, err = CurveFromStarlark(fn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid number of parameters")
}

func TestCurveFromStarlark_NonFunctionNonString(t *testing.T) {
	curve, err := CurveFromStarlark(starlark.MakeInt(42))
	assert.NoError(t, err)
	_, ok := curve.(animation.LinearCurve)
	assert.True(t, ok)
}
