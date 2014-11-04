package cache

import (
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// FileSystem is used to open an io.Reader at
// a given name string.
type FileSystem interface {
	// Open is used to return an io.Reader for a given string name.
	// If the name isn't sufficient to locate the desired io.Reader,
	// please return an os.ErrNotExist in order to be handled correctly.
	Open(name string) (io.Reader, error)
}

// Processor is used to process an io.Reader
// in to a target byte slice
type Processor interface {
	Process(io.Reader) ([]byte, error)
}

type Cache struct {
	fs    FileSystem
	pr    Processor
	name  func(string) string
	cache map[string][]byte
	mu    sync.RWMutex
}

func New(opts ...option) *Cache {
	c := &Cache{
		fs:    filesystem(os.Open),
		pr:    processor(ioutil.ReadAll),
		name:  func(s string) string { return s },
		cache: make(map[string][]byte),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Cache) Get(path string) ([]byte, bool, error) {
	// Initially obtain a read-lock
	unlocker := c.mu.RLocker()
	c.mu.RLock()
	defer func() { unlocker.Unlock() }()

	if data, ok := c.cache[path]; ok {
		return data, true, nil
	}

	// switch read-lock for write-lock
	c.mu.RUnlock()
	c.mu.Lock()
	unlocker = &c.mu

	name := c.name(path)
	rd, err := c.fs.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	data, err := c.pr.Process(rd)
	if err != nil {
		return nil, true, err
	}

	// cache the data using path
	c.cache[path] = data
	return data, true, nil
}

// filesystem is used to wrap os.Open in order
// to implement FileSystem.
type filesystem func(path string) (*os.File, error)

func (f filesystem) Open(p string) (io.Reader, error) { return f(p) }

// processor is used to wrap ioutil.ReadAll in order
// to implement Processor
type processor func(io.Reader) ([]byte, error)

func (p processor) Process(r io.Reader) ([]byte, error) { return p(r) }
