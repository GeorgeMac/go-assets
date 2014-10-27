package assets

type option func(*Assets)

func Root(path string) func(*Assets) {
	return func(a *Assets) {
		a.path = path
	}
}

func Js(js string) func(*Assets) {
	return func(a *Assets) {
		a.js = js
	}
}

func Css(css string) func(*Assets) {
	return func(a *Assets) {
		a.css = css
	}
}

func Images(images string) func(*Assets) {
	return func(a *Assets) {
		a.images = images
	}
}

func Cache(cache bool) func(*Assets) {
	return func(a *Assets) {
		a.cache = cache
	}
}
