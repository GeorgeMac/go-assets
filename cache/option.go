package cache

import (
	"fmt"
	"strings"
)

type option func(*Cache)

func Fs(fs FileSystem) option {
	return func(c *Cache) {
		c.fs = fs
	}
}

func Proc(pr Processor) option {
	return func(c *Cache) {
		c.pr = pr
	}
}

func Match(extension string) option {
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	return func(c *Cache) {
		c.name = func(s string) string {
			return fmt.Sprintf("%s%s", s, extension)
		}
		c.strip = func(s string) string {
			if strings.HasSuffix(s, extension) {
				return strings.TrimSuffix(s, extension)
			}
			return s
		}
	}
}

func PreProc(globs ...string) option {
	return func(c *Cache) {
		c.globs = globs
	}
}
