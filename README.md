`assets` package
===============

#### Usage
```go
// New asset handler for requests to '/assets/*' to serve files in './vendor/*'
asset := assets.New('/assets', assets.Dir('/vendor'))

// Use `Assets.Pattern()` for pattern arg to `http.Handle`
http.Handle(asset.Pattern(), asset)
```

#### `assets.Assets` Construction
```go
// New asset handler for pattern, with provided options
assets.New(pattern string, options ...option)
```

#### Options
```go
// Root directory of assets folder
assets.Dir(dir string)
// Cache interface
assets.SetCache(cache Cache)
```
#### Cache package

##### usage

```go
// simple in memory serving of files from disk
c := cache.New()

// in memory cache, serving files from disk and processing them with processor
c := cache.New(cache.Proc(processor))

// same as before, but matching and stripping the file extension ".md"
c := cache.New(cache.Proc(processor), cache.Match(".md"))

// cache which uses an in memory filesystem, instead of disk
fs := &InMemoryFs{
    Data: map[string][]byte{}
}
c := cache.New(cache.Fs(fs))
```

##### example
The following is an example of using `go-assets/cache` package implementation of `assets.Cache` interface.

The following is taken from the [go-assets-processors/markdown](https://github.com/GeorgeMac/go-assets-processors) package.
It is an example markdown -> html vendoring server using the go-assets `cache`.

```go
package main

import (
	"net/http"

	"github.com/GeorgeMac/go-assets"
	"github.com/GeorgeMac/go-assets-processors/markdown"
	"github.com/GeorgeMac/go-assets/cache"
)

func main() {
	md := markdown.New()
	c := cache.New(cache.Proc(md), cache.Match(".md"))
	a := assets.New("/", assets.Dir("markdown"), assets.SetCache(c))
	http.Handle(a.Pattern(), a)
	http.ListenAndServe(":8000", nil)
}
```
