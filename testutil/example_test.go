package testutil_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
	"tidbyt.dev/pixlet/testutil"
)

// TestNewTestApplet demonstrates using testutil helpers for cleaner tests
func TestNewTestApplet(t *testing.T) {
	src := `load("render.star", "render")
def main():
    return render.Root(child=render.Text("Hello"))
`
	
	app := testutil.NewTestApplet(t, "example.star", src)
	assert.NotNil(t, app)
}

// TestTableDrivenTest demonstrates table-driven tests using testutil
func TestTableDrivenTest(t *testing.T) {
	tests := []testutil.TableTest{
		{
			Name:     "valid input",
			Input:    "test",
			Expected: "result",
			WantErr:  false,
		},
		{
			Name:     "invalid input",
			Input:    "",
			Expected: "",
			WantErr:  true,
			ErrMsg:   "input cannot be empty",
		},
	}
	
	testutil.RunTableTests(t, tests, func(t *testing.T, tt testutil.TableTest) {
		// Your test logic here
		_ = tt
	})
}

// TestRunApplet demonstrates running an applet with testutil
func TestRunApplet(t *testing.T) {
	src := `load("render.star", "render")
def main():
    return render.Root(child=render.Box())
`
	
	app := testutil.NewTestApplet(t, "test.star", src)
	roots := testutil.RunApplet(t, app)
	
	require.Len(t, roots, 1)
	assert.NotNil(t, roots[0])
}

// TestRunAppletWithConfig demonstrates running an applet with config
func TestRunAppletWithConfig(t *testing.T) {
	src := `load("render.star", "render")
def main(config):
    name = config.get("name") or "World"
    return render.Root(child=render.Text("Hello, " + name))
`
	
	app := testutil.NewTestApplet(t, "test.star", src)
	config := map[string]string{"name": "Test"}
	roots := testutil.RunAppletWithConfig(t, app, config)
	
	require.Len(t, roots, 1)
}

// TestStarlarkHelpers demonstrates using Starlark value helpers
func TestStarlarkHelpers(t *testing.T) {
	src := `def test():
    return "hello"
`
	
	globals := testutil.MustParseStarlark(t, src)
	testFunc, ok := globals["test"].(*starlark.Function)
	require.True(t, ok)
	
	thread := &starlark.Thread{}
	result, err := starlark.Call(thread, testFunc, nil, nil)
	require.NoError(t, err)
	
	str := testutil.MustBeStarlarkString(t, result)
	assert.Equal(t, "hello", str)
}

// TestWithContext demonstrates using context in tests
func TestWithContext(t *testing.T) {
	src := `load("render.star", "render")
def main():
    return render.Root(child=render.Box())
`
	
	app := testutil.NewTestApplet(t, "test.star", src)
	ctx := context.Background()
	
	roots, err := app.Run(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, roots)
}

