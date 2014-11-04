package cache

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var (
	errstring string      = "Expected '%s', Got '%s'\n"
	fs        *InMemoryFs = &InMemoryFs{
		Data: map[string][]byte{
			"/test/one": []byte("here be some dataz"),
		},
	}
)

func Test_Cache_New(t *testing.T) {
	c := New()

	if c.name == nil {
		t.Errorf(errstring, "c.name to not be nil", c.name)
	}

	if c.cache == nil {
		t.Errorf(errstring, "c.cache to not be nil", c.cache)
	}

	proc := &DummyProc{}
	c = New(Fs(fs), Proc(proc))

	if !reflect.DeepEqual(fs, c.fs) {
		t.Errorf(errstring, fs, c.fs)
	}

	if !reflect.DeepEqual(proc, c.pr) {
		t.Errorf(errstring, proc, c.pr)
	}
}

func Test_Cache_Get(t *testing.T) {
	c := New(Fs(fs))

	// Get `/test/one`
	data, ok, err := c.Get("/test/one")
	// err should be nil
	if err != nil {
		t.Fatalf(errstring, "nil", err)
	}

	// ok should be true
	if !ok {
		t.Fatalf(errstring, "true", ok)
	}

	// data should be as expected
	if string(data) != "here be some dataz" {
		t.Errorf(errstring, "here be some dataz", string(data))
	}

	data, ok, err = c.Get("/test/two")
	// err should be nil
	if err != nil {
		t.Fatalf(errstring, "nil", err)
	}

	// ok should be false
	if ok {
		t.Errorf(errstring, "false", ok)
	}

	// data slice should be empty
	if len(data) > 0 {
		t.Errorf(errstring, "length of data should be 0", len(data))
	}

	proc := &DummyProc{err: errors.New("processing error")}
	c = New(Fs(fs), Proc(proc))

	data, ok, err = c.Get("/test/one")
	// err should be as expected
	if !reflect.DeepEqual(err, errors.New("processing error")) {
		t.Fatalf(errstring, errors.New("processing error"), err)
	}

	// ok should be true (as it exists)
	if !ok {
		t.Errorf(errstring, "true", ok)
	}

	// data slice should be empty
	if len(data) > 0 {
		t.Errorf(errstring, "length of data should be 0", len(data))
	}

	c = New(Fs(filesystem(func(name string) (*os.File, error) {
		return nil, errors.New("open error")
	})))

	data, ok, err = c.Get("/test/one")
	// err should be as expected
	if !reflect.DeepEqual(err, errors.New("open error")) {
		t.Fatalf(errstring, errors.New("open error"), err)
	}

	// ok should be false
	if ok {
		t.Errorf(errstring, "false", ok)
	}

	// data slice should be empty
	if len(data) > 0 {
		t.Errorf(errstring, "length of data should be 0", len(data))
	}

}

func Test_Cache_Get_OnFilesystem(t *testing.T) {
	fi, err := ioutil.TempFile("", "")
	if err != nil {
		// shouldn't be a problem
		t.Fatalf("%v\n", err)
	}

	// some initial file populating
	name := fi.Name()
	if _, err := fi.Write([]byte("dataz")); err != nil {
		// shouldn't be a problem
		t.Fatalf("%v\n", err)
	}

	c := New()
	// Get `/test/one`
	data, ok, err := c.Get(name)
	// err should be nil
	if err != nil {
		t.Fatalf(errstring, "nil", err)
	}

	// ok should be true
	if !ok {
		t.Fatalf(errstring, "true", ok)
	}

	// data should be as expected
	if string(data) != "dataz" {
		t.Errorf(errstring, "dataz", string(data))
	}
}

func Test_Cache_Get_UsesCache(t *testing.T) {
	var called bool
	fsn := testfs(func(name string) (io.Reader, error) {
		called = true
		return fs.Open(name)
	})
	c := New(Fs(fsn))

	// call get for `/test/one`
	data, ok, err := c.Get("/test/one")
	// err should be nil
	if err != nil {
		t.Errorf(errstring, "nil", err)
	}

	// ok should be true
	if !ok {
		t.Errorf(errstring, "true", ok)
	}

	// data should be as expected
	if string(data) != "here be some dataz" {
		t.Errorf(errstring, "here be some dataz", string(data))
	}

	// expect called to be true
	if !called {
		t.Fatalf(errstring, "true", called)
	}

	// reset called to false
	called = false

	// call get for `/test/one`
	data, ok, err = c.Get("/test/one")
	// err should be nil
	if err != nil {
		t.Errorf(errstring, "nil", err)
	}

	// ok should be true
	if !ok {
		t.Errorf(errstring, "true", ok)
	}

	// data should be as expected
	if string(data) != "here be some dataz" {
		t.Errorf(errstring, "here be some dataz", string(data))
	}

	// expected called to be false, due to cache hit
	if called {
		t.Fatalf(errstring, "false", called)
	}
}

func Test_Cache_Get_Matches(t *testing.T) {
	// FileSystem with fake markdown in it
	fs := &InMemoryFs{
		Data: map[string][]byte{
			"/path/to/markdown.html.md": []byte("some lovely markdown"),
		},
	}

	// new cache matching file paths with `.md` extension
	c := New(Fs(fs), Match(".md"))

	// expect c.name to not be nil
	if c.name == nil {
		t.Fatalf(errstring, "function to be present", c.name)
	}

	// expect c.name to add extension
	if c.name("somefile") != "somefile.md" {
		t.Fatalf(errstring, "somefile.md", c.name("somefile"))
	}

	// ext without `.` prefix
	c = New(Fs(fs), Match("md"))

	// expect c.name to add extension with prefixed `.`
	if c.name("otherfile") != "otherfile.md" {
		t.Fatalf(errstring, "otherfile.md", c.name("otherfile"))
	}

	// Get `/path/to/markdown.html.md` which exists, but not with matches specifier.
	data, ok, err := c.Get("/path/to/markdown.html.md")
	// err should be nil
	if err != nil {
		t.Fatalf(errstring, "nil", err)
	}

	// ok should be false
	if ok {
		t.Fatalf(errstring, "false", ok)
	}

	// data slice should be empty
	if len(data) > 0 {
		t.Errorf(errstring, "length of data should be 0", len(data))
	}

	// Get `/path/to/markdown.html` which exists with matchers specifier.
	data, ok, err = c.Get("/path/to/markdown.html")
	// err should be nil
	if err != nil {
		t.Fatalf(errstring, "nil", err)
	}

	// ok should be true
	if !ok {
		t.Fatalf(errstring, "true", ok)
	}

	// data should be as expected
	if string(data) != "some lovely markdown" {
		t.Errorf(errstring, "here be some dataz", string(data))
	}

}

// testfs implements FileSystem and delegates
// to wrapped function of the same type
type testfs func(name string) (io.Reader, error)

// Open delegates to function receiver
func (t testfs) Open(name string) (io.Reader, error) {
	return t(name)
}

// DummyProc implements cache.Processor
type DummyProc struct {
	err error
}

// Process returns an empty slice of bytes and the error within the struct
func (d *DummyProc) Process(r io.Reader) ([]byte, error) { return []byte{}, d.err }
