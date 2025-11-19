package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/render"
)

// NewTestApplet creates a new applet from Starlark source for testing.
func NewTestApplet(t *testing.T, name, src string) *runtime.Applet {
	t.Helper()
	app, err := runtime.NewApplet(name, []byte(src))
	require.NoError(t, err)
	return app
}

// RunApplet runs an applet and returns the roots, failing the test on error.
func RunApplet(t *testing.T, app *runtime.Applet) []render.Root {
	t.Helper()
	ctx := context.Background()
	roots, err := app.Run(ctx)
	require.NoError(t, err)
	return roots
}

// RunAppletWithConfig runs an applet with config and returns the roots.
func RunAppletWithConfig(t *testing.T, app *runtime.Applet, config map[string]string) []render.Root {
	t.Helper()
	ctx := context.Background()
	roots, err := app.RunWithConfig(ctx, config)
	require.NoError(t, err)
	return roots
}

// MustBeStarlarkString asserts that a value is a Starlark string and returns it.
func MustBeStarlarkString(t *testing.T, val starlark.Value) string {
	t.Helper()
	str, ok := starlark.AsString(val)
	require.True(t, ok, "expected string, got %T", val)
	return str
}

// MustBeStarlarkInt asserts that a value is a Starlark int and returns it.
func MustBeStarlarkInt(t *testing.T, val starlark.Value) int64 {
	t.Helper()
	i, ok := val.(starlark.Int)
	require.True(t, ok, "expected int, got %T", val)
	result, ok := i.Int64()
	require.True(t, ok, "int too large for int64")
	return result
}

// CompareWithGolden compares an image with a golden file.
// If UPDATE_GOLDEN is set, it updates the golden file instead.
func CompareWithGolden(t *testing.T, img image.Image, goldenPath string) {
	t.Helper()
	
	if os.Getenv("UPDATE_GOLDEN") != "" {
		// Update golden file
		require.NoError(t, os.MkdirAll(filepath.Dir(goldenPath), 0755))
		f, err := os.Create(goldenPath)
		require.NoError(t, err)
		defer f.Close()
		require.NoError(t, png.Encode(f, img))
		t.Logf("Updated golden file: %s", goldenPath)
		return
	}
	
	// Compare with golden file
	f, err := os.Open(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("Golden file does not exist: %s. Set UPDATE_GOLDEN=1 to create it.", goldenPath)
		}
		require.NoError(t, err)
	}
	defer f.Close()
	
	expected, _, err := image.Decode(f)
	require.NoError(t, err)
	
	bounds := img.Bounds()
	expectedBounds := expected.Bounds()
	
	if !bounds.Eq(expectedBounds) {
		t.Errorf("Image size mismatch: got %v, expected %v", bounds, expectedBounds)
		return
	}
	
	// Compare pixel by pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gotColor := img.At(x, y)
			expectedColor := expected.At(x, y)
			if gotColor != expectedColor {
				t.Errorf("Color mismatch at (%d, %d): got %v, expected %v", x, y, gotColor, expectedColor)
				return
			}
		}
	}
}

// LoadJSONFixture loads a JSON fixture file and unmarshals it into v.
func LoadJSONFixture(t *testing.T, path string, v interface{}) {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, v))
}

// SaveJSONFixture saves v as JSON to a fixture file.
func SaveJSONFixture(t *testing.T, path string, v interface{}) {
	t.Helper()
	data, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, data, 0644))
}

// TableTest represents a single test case in a table-driven test.
type TableTest struct {
	Name     string
	Input    interface{}
	Expected interface{}
	WantErr  bool
	ErrMsg   string // Optional: specific error message to check
}

// RunTableTests runs a table of tests using the provided test function.
func RunTableTests(t *testing.T, tests []TableTest, fn func(t *testing.T, test TableTest)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			fn(t, tt)
		})
	}
}

// AssertStarlarkEqual asserts that two Starlark values are equal.
func AssertStarlarkEqual(t *testing.T, expected, actual starlark.Value, msgAndArgs ...interface{}) {
	t.Helper()
	equal, err := starlark.Equal(expected, actual)
	if err != nil {
		t.Errorf("Error comparing Starlark values: %v", err)
		return
	}
	if !equal {
		t.Errorf("Starlark values not equal: expected %v, got %v", expected, actual)
		if len(msgAndArgs) > 0 {
			t.Logf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		}
	}
}

// CreateTestFS creates a test filesystem with the given files.
func CreateTestFS(files map[string]string) map[string][]byte {
	result := make(map[string][]byte)
	for path, content := range files {
		result[path] = []byte(content)
	}
	return result
}

// MustParseStarlark parses Starlark source and returns the globals, failing on error.
func MustParseStarlark(t *testing.T, src string) starlark.StringDict {
	t.Helper()
	thread := &starlark.Thread{}
	globals, err := starlark.ExecFile(thread, "test.star", src, nil)
	require.NoError(t, err)
	return globals
}

// GetTestDataPath returns the path to a test data file.
func GetTestDataPath(t *testing.T, relPath string) string {
	t.Helper()
	// Look for testdata in current package and parent directories
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	dir := wd
	for {
		testdataPath := filepath.Join(dir, "testdata", relPath)
		if _, err := os.Stat(testdataPath); err == nil {
			return testdataPath
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	// Fallback: return path relative to current directory
	return filepath.Join("testdata", relPath)
}

// AssertImageSize asserts that an image has the expected size.
func AssertImageSize(t *testing.T, img image.Image, expectedWidth, expectedHeight int) {
	t.Helper()
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	assert.Equal(t, expectedWidth, width, "image width")
	assert.Equal(t, expectedHeight, height, "image height")
}

// ErrorContains checks if an error message contains a substring.
func ErrorContains(t *testing.T, err error, substr string) bool {
	t.Helper()
	if err == nil {
		return false
	}
	return assert.Contains(t, err.Error(), substr, "error message should contain %q", substr)
}

// MustReadFile reads a file and returns its contents, failing on error.
func MustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file: %s", path)
	return data
}

// TempFile creates a temporary file with the given content and returns its path.
// The file is automatically cleaned up after the test.
func TempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "pixlet-test-*.tmp")
	require.NoError(t, err)
	defer f.Close()
	
	_, err = f.WriteString(content)
	require.NoError(t, err)
	
	t.Cleanup(func() {
		os.Remove(f.Name())
	})
	
	return f.Name()
}

// AssertNoError is a convenience wrapper around require.NoError with better messaging.
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if len(msgAndArgs) > 0 {
		require.NoError(t, err, msgAndArgs...)
	} else {
		require.NoError(t, err)
	}
}

// FormatTestName formats a test name with parameters for better readability.
func FormatTestName(base string, params map[string]interface{}) string {
	if len(params) == 0 {
		return base
	}
	
	result := base + "_"
	first := true
	for k, v := range params {
		if !first {
			result += "_"
		}
		result += fmt.Sprintf("%s=%v", k, v)
		first = false
	}
	return result
}

