`assets` package
===============

### API v0.2

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
#### `cache.Cache` Construction `(implements assets.Cache)`
The following is an example of using `go-assets/cache` package implementation of `assets.Cache` interface. 

Note that `type LessProcessor struct` does not exist. It is an example of an implementation of the `cache.Processor` interface, that is yet to exist.

```go
lessCache := cache.New(cache.Proc(&LessProcessor{}), 
	cache.Match(".less"))

// example:
// GET /assets/default.css
// reads /vendor/less/default.css.less on disk
// writes result of Processor.Processor as response
// and caches it.
lessAssets := assets.New("/assets", 
	assets.Dir("/vendor/less"), 
	assets.Cache(lessCache))
	
http.Handle(lessAssets.Pattern(), lessAssets)
```




