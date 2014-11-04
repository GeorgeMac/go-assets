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
	return func(c *Cache) {
		c.name = func(s string) string {
			if strings.HasPrefix(extension, ".") {
				return s + extension
			}
			return fmt.Sprintf("%s.%s", s, extension)
		}
	}
}
