package assets

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

	if a.pattern != "/assets/" {
		t.Fatalf(errstring, "/assets/", a.pattern)
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

	a = New("/")

	if a.Pattern() != "/" {
		t.Fatalf(errstring, "/", a.Pattern())
	}
}

func Test_Assets_ServeHTTP_BadCache(t *testing.T) {
	// set up logging
	log := logger()

	// [test] InternalServerError on cache error
	// new Assets struct, with bad cache
	a := New("/assets", SetCache(tcache(func(_ string) ([]byte, bool, error) {
		return nil, false, errors.New("cache error")
	})))

	// build bad request
	req := request("/", nil)
	// get a new recorder
	resp := httptest.NewRecorder()

	// run request handler
	a.ServeHTTP(resp, req)

	// Check response code is as expected
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf(errstring, http.StatusInternalServerError, resp.Code)
	}

	// check log message for expected error
	msg := log.String()
	if msg != "[ERROR] cache error\n" {
		t.Errorf(errstring, "[ERROR] cache error\n", msg)
	}
}

func Test_Assets_ServeHTTP_NotFound(t *testing.T) {
	// set up logging
	log := logger()

	// [test] FileNotFound on cache error
	// new Assets struct, with empty cache
	a := New("/assets", SetCache(tcache(func(_ string) ([]byte, bool, error) {
		return nil, false, nil
	})))

	// build bad request
	req := request("/assets/stuff", nil)
	// get a new recorder
	resp := httptest.NewRecorder()

	// run request handler
	a.ServeHTTP(resp, req)

	// Check response code is as expected
	if resp.Code != http.StatusNotFound {
		t.Fatalf(errstring, http.StatusNotFound, resp.Code)
	}

	// check log message for expected error
	msg := log.String()
	if msg != "[INFO] not found /assets/stuff\n" {
		t.Errorf(errstring, "[INFO] not found /assets/stuff\n", msg)
	}
}

func Test_Assets_ServeHTTP_OK(t *testing.T) {
	// [test] StatusOK, file is delivered and Content-Type is correct

	// set up in memory filesystem for cache
	fs := &cache.InMemoryFs{
		Data: map[string][]byte{
			"/vendor/one":      []byte("Hello, World!"),
			"/vendor/two.html": []byte("<html><body>Hello, World!</body></html>"),
			"/vendor/three.js": []byte("console.Log(\"Hello, World!\")"),
		},
	}

	// construct a new cache
	c := cache.New(cache.Fs(fs))

	// new Assets struct, with empty cache
	a := New("/assets", Dir("/vendor"), SetCache(c))

	type TestCase struct {
		path, ctype string
		body        []byte
	}

	cases := []TestCase{
		TestCase{
			path:  "/assets/one",
			body:  fs.Data["/vendor/one"],
			ctype: "text/plain",
		},
		TestCase{
			path:  "/assets/two.html",
			body:  fs.Data["/vendor/two.html"],
			ctype: "text/html; charset=utf-8",
		},
		TestCase{
			path:  "/assets/three.js",
			body:  fs.Data["/vendor/three.js"],
			ctype: "application/javascript",
		},
	}

	for _, cc := range cases {
		// build bad request
		req := request(cc.path, nil)
		// get a new recorder
		resp := httptest.NewRecorder()

		// run request handler
		a.ServeHTTP(resp, req)

		// Check response code is as expected
		if resp.Code != http.StatusOK {
			t.Fatalf(errstring, http.StatusOK, resp.Code)
		}

		// check log message for expected error
		body := resp.Body.Bytes()
		if !reflect.DeepEqual(body, cc.body) {
			t.Errorf(errstring, string(cc.body), string(body))
		}

		// check Content-Type header
		ctype := resp.HeaderMap["Content-Type"][0]
		if ctype != cc.ctype {
			t.Errorf(errstring, cc.ctype, ctype)
		}
	}
}

// set up log package with a buffer to call String() on
func logger() fmt.Stringer {
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	log.SetFlags(0)
	return buf
}

// build a request to the specified path, with the provided body
func request(path string, body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com%s", path), body)
	if err != nil {
		panic(err)
	}
	return req
}

// tcache implements Cache
// used to construct Cache type from anonymous function
type tcache func(string) ([]byte, bool, error)

func (c tcache) Get(path string) ([]byte, bool, error) {
	return c(path)
}
