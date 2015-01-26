package cache

import (
	"bytes"
	"io"
	"os"
	"path"
)

type InMemoryFs struct {
	Data map[string][]byte
}

func (fs *InMemoryFs) Open(name string) (io.Reader, error) {
	data, ok := fs.Data[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	return bytes.NewBuffer(data), nil
}

func (fs *InMemoryFs) Glob(glob string) ([]string, error) {
	matches := []string{}
	for k := range fs.Data {
		ok, err := path.Match(glob, k)
		if err != nil {
			return matches, err
		}
		if ok {
			matches = append(matches, k)
		}
	}
	return matches, nil
}
