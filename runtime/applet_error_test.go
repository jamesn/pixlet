package runtime

import (
	"context"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppletContextCancellation(t *testing.T) {
	// Test that context cancellation is handled properly
	// Use a long sleep to test cancellation
	src := `load("render.star", "render")
load("time.star", "time")
def main():
    # Sleep for a long time to allow cancellation
    time.sleep(1.0)
    return render.Root(child=render.Box())
`
	app, err := NewApplet("test.star", []byte(src))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancel context immediately
	cancel()

	_, err = app.Run(ctx)
	assert.Error(t, err)
}

func TestAppletTimeout(t *testing.T) {
	// Test that timeout is handled properly
	src := `load("render.star", "render")
load("time.star", "time")
def main():
    # Sleep for a long time
    time.sleep(1.0)
    return render.Root(child=render.Box())
`
	app, err := NewApplet("test.star", []byte(src))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = app.Run(ctx)
	assert.Error(t, err)
}

func TestExtractRootsNone(t *testing.T) {
	// Test ExtractRoots with None value
	src := `def main():
    return None
`
	app, err := NewApplet("test.star", []byte(src))
	require.NoError(t, err)

	roots, err := app.Run(context.Background())
	require.NoError(t, err)
	assert.Empty(t, roots, "None should return empty roots")
}

func TestExtractRootsInvalidType(t *testing.T) {
	// Test ExtractRoots with invalid return type
	src := `load("render.star", "render")
def main():
    return "not a root"
`
	app, err := NewApplet("test.star", []byte(src))
	require.NoError(t, err)

	_, err = app.Run(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Root", "error should mention Root")
}

func TestExtractRootsListWithInvalid(t *testing.T) {
	// Test ExtractRoots with list containing invalid items
	src := `load("render.star", "render")
def main():
    return [render.Root(child=render.Box()), "invalid"]
`
	app, err := NewApplet("test.star", []byte(src))
	require.NoError(t, err)

	_, err = app.Run(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "index", "error should mention index")
}

func TestMultipleMainFunctions(t *testing.T) {
	// Test that multiple main functions are detected
	fsys := fstest.MapFS{
		"a.star": &fstest.MapFile{
			Data: []byte(`def main():
    return None
`),
		},
		"b.star": &fstest.MapFile{
			Data: []byte(`def main():
    return None
`),
		},
	}

	_, err := NewAppletFromFS("test", fsys)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple files with a main()", "error should mention multiple main functions")
}

func TestSchemaHandlerNotFound(t *testing.T) {
	// Test that missing schema handlers return errors
	// Create an applet with a schema but no handlers
	src := `load("render.star", "render")
load("schema.star", "schema")
def get_schema():
    return schema.Schema(
        version = "1",
    )
def main():
    return render.Root(child=render.Box())
`
	app, err := NewApplet("test.star", []byte(src))
	require.NoError(t, err)

	// Try to call a handler that doesn't exist
	_, err = app.CallSchemaHandler(context.Background(), "nonexistent", "param")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no exported handler", "error should mention handler not found")
}

