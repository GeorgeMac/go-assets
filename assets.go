package assets

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/GeorgeMac/go-assets/cache"
)

type Cache interface {
	Get(path string) ([]byte, bool, error)
}

type Assets struct {
	pattern, dir string
	cache        Cache
}

func New(pattern string, options ...option) *Assets {
	a := &Assets{
		pattern: pattern,
		dir:     pattern,
		cache:   cache.New(),
	}

	for _, option := range options {
		option(a)
	}

	if a.pattern[len(a.pattern)-1] != '/' {
		a.pattern += "/"
	}

	return a
}

func (a *Assets) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	subpath, err := filepath.Rel(a.pattern, req.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] %v\n", err)
		return
	}

	dpath := filepath.Join(a.dir, subpath)
	data, ok, err := a.cache.Get(dpath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] %v\n", err)
		return
	}

	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		log.Printf("[INFO] not found %s\n", dpath)
		return
	}

	if _, err := rw.Write(data); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] %v\n", err)
		log.Printf("[INFO] %s\n", string(data))
	}

	// derive Content-Type header
	ext := filepath.Ext(subpath)
	ct := mime.TypeByExtension(ext)
	if ext == "" {
		ct = "text/plain"
	}
	rw.Header().Set("Content-Type", ct)
}

func (a *Assets) Pattern() string { return a.pattern }
