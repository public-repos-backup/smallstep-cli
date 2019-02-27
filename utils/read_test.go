package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper function for setting os.Stdin for mocking in tests.
func setStdin(new *os.File) (cleanup func()) {
	old := _osStdin
	_osStdin = new
	return func() { _osStdin = old }
}

func TestFileExists(t *testing.T) {
	content := []byte("my file content")
	f, cleanup := newFile(t, content)
	defer cleanup()

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ok", args{f.Name()}, true},
		{"nok", args{f.Name() + ".foo"}, false},
		{"empty", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.path); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	content := []byte("my file content")
	f, cleanup := newFile(t, content)
	defer cleanup()

	b, err := ReadFile(f.Name())
	require.NoError(t, err)
	require.True(t, bytes.Equal(content, b), "expected %s to equal %s", b, content)
}

func TestReadFileStdin(t *testing.T) {
	content := []byte("my file content")
	mockStdin, cleanup := newFile(t, content)
	defer cleanup()
	defer setStdin(mockStdin)()

	b, err := ReadFile(stdinFilename)
	require.NoError(t, err)
	require.True(t, bytes.Equal(content, b), "expected %s to equal %s", b, content)
}

func TestReadPasswordFromFile(t *testing.T) {
	content := []byte("my-password-on-file\n")
	f, cleanup := newFile(t, content)
	defer cleanup()

	b, err := ReadPasswordFromFile(f.Name())
	require.NoError(t, err)
	require.True(t, bytes.Equal([]byte("my-password-on-file"), b), "expected %s to equal %s", b, content)
}

func TestStringReadPasswordFromFile(t *testing.T) {
	content := []byte("my-password-on-file\n")
	f, cleanup := newFile(t, content)
	defer cleanup()

	s, err := ReadStringPasswordFromFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, "my-password-on-file", s, "expected %s to equal %s", s, content)
}

// Returns a temp file and a cleanup function to delete it.
func newFile(t *testing.T, data []byte) (file *os.File, cleanup func()) {
	f, err := ioutil.TempFile("" /* dir */, "utils-read-test")
	require.NoError(t, err)
	// write to temp file and reset read cursor to beginning of file
	_, err = f.Write(data)
	require.NoError(t, err)
	_, err = f.Seek(0, io.SeekStart)
	require.NoError(t, err)
	return f, func() { os.Remove(f.Name()) }
}
