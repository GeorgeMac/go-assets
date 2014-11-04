package assets

import (
	"reflect"
	"testing"

	"github.com/GeorgeMac/go-assets/cache"
)

var errstring string = "Expected '%s', Got '%s'\n"

func Test_Assets_New(t *testing.T) {
	a := New("/")

	if a.dir != "/" {
		t.Fatalf(errstring, "/", a.dir)
	}

	if a.pattern != "/" {
		t.Fatalf(errstring, "/", a.pattern)
	}

	if a.cache == cache.New() {
		t.Fatalf(errstring, cache.New(), a.cache)
	}

	c := cache.New(cache.Match("md"))
	a = New("/assets", Dir("/vendor"), SetCache(c))

	if a.pattern != "/assets" {
		t.Fatalf(errstring, "/assets", a.pattern)
	}

	if a.dir != "/vendor" {
		t.Fatalf(errstring, "/vendor", a.dir)
	}

	if !reflect.DeepEqual(a.cache, c) {
		t.Fatalf(errstring, c, a.cache)
	}
}

func Test_Assets_Pattern(t *testing.T) {
	a := New("/assets", Dir("/vendor"))

	// check pattern is as expected with trailing slash
	if a.Pattern() != "/assets/" {
		t.Fatalf(errstring, "/assets/", a.Pattern())
	}
}

func Test_Assets_ServeHttp(t *testing.T) {

}
