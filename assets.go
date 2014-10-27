package assets

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"sync"
)

type Assets struct {
	path, js, css, images string
	cache                 bool
	acache                map[string][]byte
	mu                    sync.RWMutex
}

func New(options ...option) *Assets {
	a := &Assets{
		path:   "assets",
		js:     "js",
		css:    "css",
		images: "images",
		acache: map[string][]byte{},
	}

	for _, opt := range options {
		opt(a)
	}

	return a
}

func (a *Assets) Path() string {
	return path.Join("/", a.path) + "/"
}

func (a *Assets) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			_, asset := path.Split(req.URL.Path)
			typ, fpath := path.Split(asset)
			switch typ {
			case a.js:
				w.Header().Set("Content-Type", "application/javascript")
			case a.css:
				w.Header().Set("Content-Type", "text/css")
			case a.images:
				ext := path.Ext(fpath)
				w.Header().Set("Content-Type", mime.TypeByExtension(ext))
			default:
				w.Header().Set("Content-Type", "text/plain")
			}

			// derive resouce path
			resource := path.Join("./", a.path, typ, fpath)
			if a.cache {
				// check cache
				a.mu.RLock()
				log.Printf("[DEBUG] Checking cache for %s\n", resource)
				if data, ok := a.acache[resource]; ok {
					if _, err := w.Write(data); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
					a.mu.RUnlock()
					log.Println("[DEBUG] Cache hit!")
					return
				}
				log.Println("[DEBUG] Cache miss!")
				a.mu.RUnlock()

				a.mu.Lock()
				defer a.mu.Unlock()

				data, err := ioutil.ReadFile(resource)
				if err != nil {
					if os.IsNotExist(err) {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if _, err := w.Write(data); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				log.Printf("[DEBUG] Caching resource %s\n", resource)
				a.acache[resource] = data
				return
			}

			// no caching
			log.Println("[DEBUG] Caching disabled")
			fi, err := os.Open(resource)
			if err != nil {
				if os.IsNotExist(err) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(w, fi); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
