package fuel

import (
	"fmt"
	"reflect"
	"strings"
)

const ignored = "-"

type Fixture struct {
	Parent *Fixture

	// base url
	Prefix  string
	Root    string
	Route   string
	Version string

	// middleware
	Middleware []string

	// stub
	Stub string

	// caching
	Cache string
	TTL   string
}

func newFixture(tag reflect.StructTag) Fixture {

	read := func(tag reflect.StructTag, keys ...string) string {
		var val string
		for _, k := range keys {
			val = tag.Get(k)
			if val != "" {
				return val
			}
		}
		return ""
	}

	return Fixture{
		Prefix:  read(tag, "prefix", "pre"),
		Root:    read(tag, "root"),
		Route:   read(tag, "route"),
		Version: read(tag, "version", "ver", "v"),
		Cache:   read(tag, "cache"),
		TTL:     read(tag, "ttl"),
		Stub:    read(tag, "stub"),
		Middleware: func() []string {
			m := []string{}
			list := strings.Split(read(tag, "middle", "middleware"), ",")
			for _, l := range list {
				midw := strings.TrimSpace(l)
				if midw != "" {
					m = append(m, midw)
				}
			}
			if len(m) == 0 {
				return nil
			}
			return m
		}(),
	}
}

func (f Fixture) getPrefix() string {
	value := f.Prefix
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getPrefix()
	}

	return value
}

func (f Fixture) getRoot() string {
	value := f.Root
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getRoot()
	}

	return value
}

func (f Fixture) getRoute() string {
	value := f.Route
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getRoute()
	}

	return value
}

func (f Fixture) getURL() string {

	var urlOut string

	var prefix, version, root, url = f.getPrefix(), f.getVersion(), f.getRoot(), f.getRoute()

	if prefix != "" {
		prefix = "/" + prefix
	}
	if version != "" {
		version = "/v" + version
	}
	if root != "" {
		root = "/" + root
	}
	if url != "" {
		url = "/" + url
	}

	if VersionAfterPrefix == true {
		urlOut = fmt.Sprintf("%s%s%s%s", prefix, version, root, url)
	} else {
		urlOut = fmt.Sprintf("%s%s%s%s", prefix, root, url, version)
	}

	return cleanMultSlash(urlOut)
}

func (f Fixture) getVersion() string {
	value := f.Version
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getVersion()
	}

	return value
}

func (f Fixture) getMiddleware() []string {
	value := f.Middleware

	if (value == nil || len(value) == 0) && f.Parent != nil {
		value = f.Parent.getMiddleware()
	}

	return value
}

func (f Fixture) getStub() string {
	value := f.Stub
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getStub()
	}

	return value
}

func (f Fixture) getCache() string {
	value := f.Cache
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getCache()
	}

	return value
}

func (f Fixture) getTTL() string {
	value := f.TTL
	if value == ignored {
		return ""
	}

	if value == "" && f.Parent != nil {
		value = f.Parent.getTTL()
	}

	return value
}
