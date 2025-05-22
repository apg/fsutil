package union

import (
	"errors"
	"io/fs"
	"path"
)

type FS struct {
	layers []fs.FS
}

// New creates a new UnionFS where files are resolved in the order given.
// The first FS to return a successful Open will be used.
func New(layers ...fs.FS) *FS {
	return &FS{layers: layers}
}

// Open implements the FS interface
func (u *FS) Open(name string) (fs.File, error) {
	cleaned := path.Clean(name)
	for _, l := range u.layers {
		file, err := l.Open(cleaned)
		if err == nil {
			return file, nil
		}

		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	return nil, fs.ErrNotExist
}

// ReadFile implements the fs.ReadFileFS interface
func (u *FS) ReadFile(name string) ([]byte, error) {
	for _, l := range u.layers {
		bs, err := fs.ReadFile(l, name)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		return bs, err
	}
	return nil, fs.ErrNotExist
}
