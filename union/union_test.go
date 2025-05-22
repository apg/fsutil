package union

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestUnionFS_OpenPrefersFirstMatch(t *testing.T) {
	fs1 := fstest.MapFS{
		"foo.txt": &fstest.MapFile{Data: []byte("from fs1")},
	}
	fs2 := fstest.MapFS{
		"foo.txt": &fstest.MapFile{Data: []byte("from fs2")},
		"bar.txt": &fstest.MapFile{Data: []byte("only in fs2")},
	}

	union := New(fs1, fs2)

	// Should prefer fs1
	data, err := fs.ReadFile(union, "foo.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "from fs1" {
		t.Errorf("expected 'from fs1', got '%s'", data)
	}

	// Should fallback to fs2
	data, err = fs.ReadFile(union, "bar.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "only in fs2" {
		t.Errorf("expected 'only in fs2', got '%s'", data)
	}
}

func TestUnionFS_MissingFile(t *testing.T) {
	union := New(fstest.MapFS{})

	_, err := fs.ReadFile(union, "missing.txt")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("expected not exist error, got: %v", err)
	}
}
