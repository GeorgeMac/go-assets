package assets

import "testing"

var expstr string = "Expected a.%s to be '%s', instead found '%s'"

func Test_AssetsNew(t *testing.T) {
	a := New()

	if a.path != "assets" {
		t.Errorf(expstr, "assets", "assets", a.path)
	}

	if a.js != "js" {
		t.Errorf(expstr, "js", "js", a.js)
	}

	if a.css != "css" {
		t.Errorf(expstr, "css", "css", a.css)
	}

	if a.images != "images" {
		t.Errorf(expstr, "images", "images", a.images)
	}

	if a.cache {
		t.Errorf(expstr, "cache", "false", a.cache)
	}

	path := a.Path()
	if path != "/assets/" {
		t.Errorf(expstr, "Path()", "/assets/", path)
	}

	a = New(Root("a"), Js("b"), Css("c"), Images("d"), Cache(true))

	if a.path != "a" {
		t.Errorf(expstr, "assets", "a", a.path)
	}

	if a.js != "b" {
		t.Errorf(expstr, "js", "b", a.js)
	}

	if a.css != "c" {
		t.Errorf(expstr, "css", "c", a.css)
	}

	if a.images != "d" {
		t.Errorf(expstr, "images", "d", a.images)
	}

	if !a.cache {
		t.Errorf(expstr, "cache", "true", a.cache)
	}

	path = a.Path()
	if path != "/a/" {
		t.Errorf(expstr, "Path()", "/a/", path)
	}

}
