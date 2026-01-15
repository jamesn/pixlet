package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
)

func TestAppletConfigGetString(t *testing.T) {
	config := AppletConfig{
		"key1": "value1",
		"key2": "value2",
	}

	thread := &starlark.Thread{}

	// Test getting existing key
	result, err := config.getString(thread, nil, starlark.Tuple{starlark.String("key1")}, nil)
	require.NoError(t, err)
	assert.Equal(t, "value1", string(result.(starlark.String)))

	// Test getting non-existent key with default
	result, err = config.getString(thread, nil, starlark.Tuple{starlark.String("missing"), starlark.String("default")}, nil)
	require.NoError(t, err)
	assert.Equal(t, "default", string(result.(starlark.String)))

	// Test getting non-existent key without default
	result, err = config.getString(thread, nil, starlark.Tuple{starlark.String("missing")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.None, result)
}

func TestAppletConfigGetBoolean(t *testing.T) {
	config := AppletConfig{
		"true_val":  "true",
		"false_val": "false",
		"one_val":   "1",
		"zero_val":  "0",
		"invalid":   "not_a_bool",
	}

	thread := &starlark.Thread{}

	// Test true values
	result, err := config.getBoolean(thread, nil, starlark.Tuple{starlark.String("true_val")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.Bool(true), result)

	result, err = config.getBoolean(thread, nil, starlark.Tuple{starlark.String("one_val")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.Bool(true), result)

	// Test false values
	result, err = config.getBoolean(thread, nil, starlark.Tuple{starlark.String("false_val")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.Bool(false), result)

	result, err = config.getBoolean(thread, nil, starlark.Tuple{starlark.String("zero_val")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.Bool(false), result)

	// Test invalid boolean
	result, err = config.getBoolean(thread, nil, starlark.Tuple{starlark.String("invalid")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.Bool(false), result)

	// Test missing key with default
	result, err = config.getBoolean(thread, nil, starlark.Tuple{starlark.String("missing"), starlark.Bool(true)}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.Bool(true), result)

	// Test missing key without default
	result, err = config.getBoolean(thread, nil, starlark.Tuple{starlark.String("missing")}, nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.None, result)
}

func TestAppletConfigAttr(t *testing.T) {
	config := AppletConfig{}

	// Test str attribute
	val, err := config.Attr("str")
	require.NoError(t, err)
	assert.NotNil(t, val)
	_, ok := val.(*starlark.Builtin)
	assert.True(t, ok)

	// Test get attribute (alias for str)
	val, err = config.Attr("get")
	require.NoError(t, err)
	assert.NotNil(t, val)
	_, ok = val.(*starlark.Builtin)
	assert.True(t, ok)

	// Test bool attribute
	val, err = config.Attr("bool")
	require.NoError(t, err)
	assert.NotNil(t, val)
	_, ok = val.(*starlark.Builtin)
	assert.True(t, ok)

	// Test invalid attribute
	val, err = config.Attr("invalid")
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestAppletConfigAttrNames(t *testing.T) {
	config := AppletConfig{}
	names := config.AttrNames()
	
	assert.Contains(t, names, "get")
	assert.Contains(t, names, "str")
	assert.Contains(t, names, "bool")
	assert.Equal(t, 3, len(names))
}

func TestAppletConfigGet(t *testing.T) {
	config := AppletConfig{
		"key1": "value1",
	}

	// Test getting existing key
	val, found, err := config.Get(starlark.String("key1"))
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value1", string(val.(starlark.String)))

	// Test getting non-existent key
	val, found, err = config.Get(starlark.String("missing"))
	require.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, "", string(val.(starlark.String)))

	// Test with non-string key
	val, found, err = config.Get(starlark.MakeInt(42))
	require.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestAppletConfigStarlarkMethods(t *testing.T) {
	config := AppletConfig{"key": "value"}

	assert.Equal(t, "AppletConfig(...)", config.String())
	assert.Equal(t, "AppletConfig", config.Type())
	assert.Equal(t, starlark.Bool(true), config.Truth())

	// Test Hash
	hash, err := config.Hash()
	require.NoError(t, err)
	assert.NotEqual(t, uint32(0), hash)

	// Test Freeze (should not panic)
	assert.NotPanics(t, func() {
		config.Freeze()
	})
}
