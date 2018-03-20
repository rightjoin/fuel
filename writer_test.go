package fuel

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {

	type Structy struct {
		FieldA string
		FieldB string
	}

	// encode an array of supported items
	r := []reflect.Value{
		reflect.ValueOf(12345),   // int
		reflect.ValueOf("54321"), // string
		reflect.ValueOf(map[string]interface{}{ // map
			"key1": "value1",
			"key2": "value2",
		}),
		reflect.ValueOf(Structy{ // struct
			FieldA: "a",
			FieldB: "b",
		}),
		reflect.ValueOf([]Structy{ // array of struct
			Structy{FieldA: "a1", FieldB: "b1"},
			Structy{FieldA: "a2", FieldB: "b2"},
		}),
	}

	//extract types
	typ := make([]reflect.Type, len(r))
	for i := range r {
		typ[i] = r[i].Type()
	}

	// write, and read
	s := serializer{}
	buf := s.write(r)
	vals := s.read(buf, typ)

	// validate values
	assert.Equal(t, int64(12345), vals[0].Int())
	assert.Equal(t, "54321", vals[1].String())

	assert.Equal(t, "value1", vals[2].Interface().(map[string]interface{})["key1"])
	assert.Equal(t, "value2", vals[2].Interface().(map[string]interface{})["key2"])

	// struct is read as map
	assert.Equal(t, "a", vals[3].Interface().(map[string]interface{})["FieldA"])

	// slice is read as []interface{}
	assert.Equal(t, "a1", vals[4].Interface().([]interface{})[0].(map[string]interface{})["FieldA"])
}
