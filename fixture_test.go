package fuel

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixtureParse(t *testing.T) {

	type AnyStruct struct {
		anyField GET `prefix:"abc" middle:"abcd, efgh , ijkl"`
	}
	var a AnyStruct

	utype := reflect.TypeOf(a)
	field, _ := utype.FieldByName("anyField")

	f := newFixture(field.Tag)

	assert.Equal(t, "abc", f.Prefix)

	assert.Equal(t, 3, len(f.Middleware))
	assert.Equal(t, "abcd", f.Middleware[0])
	assert.Equal(t, "efgh", f.Middleware[1])
	assert.Equal(t, "ijkl", f.Middleware[2])
}

func TestFixtureOverride(t *testing.T) {

	parent := Fixture{
		Prefix:  "p1",
		Root:    "r1",
		Version: "v1",
		Route:   "r1",
		Cache:   "c1",
		TTL:     "t1",
	}

	child1 := Fixture{
		Parent: &parent,
	}
	assert.Equal(t, "p1", child1.getPrefix())
	assert.Equal(t, "r1", child1.getRoot())
	assert.Equal(t, "v1", child1.getVersion())
	assert.Equal(t, "r1", child1.getRoute())
	assert.Equal(t, "c1", child1.getCache())
	assert.Equal(t, "t1", child1.getTTL())

	child2 := Fixture{
		Prefix:  "p2",
		Root:    "r2",
		Version: "v2",
		Route:   "r2",
		Cache:   "c2",
		TTL:     "t2",
		Parent:  &parent,
	}
	assert.Equal(t, "p2", child2.getPrefix())
	assert.Equal(t, "r2", child2.getRoot())
	assert.Equal(t, "v2", child2.getVersion())
	assert.Equal(t, "r2", child2.getRoute())
	assert.Equal(t, "c2", child2.getCache())
	assert.Equal(t, "t2", child2.getTTL())
}

func TestURL(t *testing.T) {
	fix := Fixture{
		Prefix:  "p1",
		Root:    "r1",
		Version: "v1",
		Route:   "r1",
	}

	VersionAfterPrefix = true
	assert.Equal(t, "/p1/vv1/r1/r1", fix.getURL())

	VersionAfterPrefix = false
	assert.Equal(t, "/p1/r1/r1/vv1", fix.getURL())
}
