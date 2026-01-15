package file

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestFileReadAll(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("hello world"),
		},
	}

	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	thread := &starlark.Thread{}

	// Test readall without mode
	result, err := f.readall(thread, nil, starlark.Tuple{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(result.(starlark.String)))

	// Test readall with text mode
	result, err = f.readall(thread, nil, starlark.Tuple{starlark.String("r")}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(result.(starlark.String)))

	// Test readall with binary mode
	result, err = f.readall(thread, nil, starlark.Tuple{starlark.String("rb")}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(result.(starlark.Bytes)))
}

func TestFileReadAllInvalidMode(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("hello"),
		},
	}

	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	thread := &starlark.Thread{}

	// Test invalid mode
	_, err := f.readall(thread, nil, starlark.Tuple{starlark.String("w")}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported mode")
}

func TestFileReadAllNonExistentFile(t *testing.T) {
	fs := fstest.MapFS{}

	f := &File{
		FS:   fs,
		Path: "nonexistent.txt",
	}

	thread := &starlark.Thread{}

	_, err := f.readall(thread, nil, starlark.Tuple{}, nil)
	assert.Error(t, err)
}

func TestFileAttr(t *testing.T) {
	fs := fstest.MapFS{}
	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	// Test path attribute
	val, err := f.Attr("path")
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", string(val.(starlark.String)))

	// Test readall attribute
	val, err = f.Attr("readall")
	assert.NoError(t, err)
	assert.NotNil(t, val)
	_, ok := val.(*starlark.Builtin)
	assert.True(t, ok)

	// Test invalid attribute
	val, err = f.Attr("invalid")
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestFileAttrNames(t *testing.T) {
	f := &File{}
	names := f.AttrNames()
	assert.Contains(t, names, "path")
	assert.Contains(t, names, "readall")
	assert.Equal(t, 2, len(names))
}

func TestFileType(t *testing.T) {
	f := &File{}
	assert.Equal(t, "File", f.Type())
	assert.Equal(t, "File(...)", f.String())
	assert.Equal(t, starlark.Bool(true), f.Truth())
}

func TestFileHash(t *testing.T) {
	f1 := &File{Path: "test.txt"}
	f2 := &File{Path: "test.txt"}
	f3 := &File{Path: "other.txt"}

	h1, err1 := f1.Hash()
	h2, err2 := f2.Hash()
	h3, err3 := f3.Hash()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	// Same files should have same hash
	assert.Equal(t, h1, h2)

	// Different files should have different hash
	assert.NotEqual(t, h1, h3)
}

func TestReaderRead(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("hello world"),
		},
	}

	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	// Test text mode reader
	reader, err := f.reader("r")
	assert.NoError(t, err)
	defer reader.Close()

	thread := &starlark.Thread{}

	// Read all
	result, err := reader.read(thread, nil, starlark.Tuple{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(result.(starlark.String)))
}

func TestReaderReadWithSize(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("hello world"),
		},
	}

	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	reader, err := f.reader("r")
	assert.NoError(t, err)
	defer reader.Close()

	thread := &starlark.Thread{}

	// Read first 5 bytes
	result, err := reader.read(thread, nil, starlark.Tuple{starlark.MakeInt(5)}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(result.(starlark.String)))
}

func TestReaderBinaryMode(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte{0x01, 0x02, 0x03},
		},
	}

	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	reader, err := f.reader("rb")
	assert.NoError(t, err)
	defer reader.Close()

	thread := &starlark.Thread{}

	result, err := reader.read(thread, nil, starlark.Tuple{}, nil)
	assert.NoError(t, err)

	bytes := result.(starlark.Bytes)
	assert.Equal(t, byte(0x01), bytes[0])
	assert.Equal(t, byte(0x02), bytes[1])
	assert.Equal(t, byte(0x03), bytes[2])
}

func TestReaderStruct(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("test"),
		},
	}

	f := &File{
		FS:   fs,
		Path: "test.txt",
	}

	reader, err := f.reader("r")
	assert.NoError(t, err)

	s := reader.Struct()
	assert.NotNil(t, s)

	// Check that read method exists
	readVal, err := s.Attr("read")
	assert.NoError(t, err)
	assert.NotNil(t, readVal)

	// Check that close method exists
	closeVal, err := s.Attr("close")
	assert.NoError(t, err)
	assert.NotNil(t, closeVal)
}
