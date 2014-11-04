package assets

type option func(a *Assets)

func Dir(dir string) option {
	return func(a *Assets) {
		a.dir = dir
	}
}

func SetCache(cache Cache) option {
	return func(a *Assets) {
		a.cache = cache
	}
}
