`assets` package
===============


Very simple asset vendoring package for `net/http` ServerMux

## Example Usage

```go
pipeline := assets.New()

// handles assets in the following way:
// /assets/js/file.js -> javascipt
// /assets/css/file.css -> css
// /assets/images/file.png -> png image
http.Handle(pipeline.Path(), pipeline.Handler())

// variadic option setting
pipeline = assets.New(assets.Root(“vendor”), assets.Js(“javascript”), assets.Css(“stylesheet”), assets.Images(“ims”), assets.Cache(true))

// handles assets in the following way (+ it caches):
// /vendor/javascript/file.js -> javascipt
// /vendor/stylesheet/file.css -> css
// /vendor/ims/file.png -> png image
http.Handle(pipeline.Path(), pipeline.Handler())
```

### Options

```go
assets.Root(pth string)
assets.Js(pth string)
assets.Css(pth string)
assets.Images(pth string)
assrts.Cache(cache bool)
```
