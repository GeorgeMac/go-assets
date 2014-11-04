package cache

import (
	"bytes"
	"io"
	"os"
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
