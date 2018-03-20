package fuel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type writer interface {
	write(r []reflect.Value) *bytes.Buffer
	read(buf *bytes.Buffer, typ []reflect.Type) []reflect.Value
}

type serializer struct {
}

func (s *serializer) write(r []reflect.Value) *bytes.Buffer {
	buf := new(bytes.Buffer)
	encd := json.NewEncoder(buf)

	for _, ri := range r {
		s.writeItem(encd, ri, typeSymbol(ri.Type()))
	}

	return buf
}

func (s *serializer) writeItem(j *json.Encoder, r reflect.Value, t string) {

	switch {
	case t == "int":
		err := j.Encode(r.Int())
		if err != nil {
			panic(err)
		}
	case t == "map":
		err := j.Encode(r.Interface().(map[string]interface{}))
		if err != nil {
			panic(err)
		}
	case t == "string":
		err := j.Encode(r.String())
		if err != nil {
			panic(err)
		}
	case t == "i:.":
		sym := typeSymbol(reflect.TypeOf(r.Interface()))
		err := j.Encode(s)
		if err != nil {
			panic(err)
		}
		s.writeItem(j, r, sym)
	case strings.HasPrefix(t, "st:"):
		err := j.Encode(r.Interface())
		if err != nil {
			panic(err)
		}
	case strings.HasPrefix(t, "sl:"):
		err := j.Encode(r.Interface())
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("cannot write/encode '%s' for caching", t))
	}
}

func (s *serializer) read(buf *bytes.Buffer, typ []reflect.Type) []reflect.Value {
	decd := json.NewDecoder(buf)
	r := make([]reflect.Value, len(typ))
	for i := range typ {
		r[i] = s.readItem(decd, typ[i])
	}
	return r
}

func (s *serializer) readItem(j *json.Decoder, typ reflect.Type) reflect.Value {
	var r reflect.Value
	t := typeSymbol(typ)
	switch {
	case t == "int":
		var i int
		err := j.Decode(&i)
		if err != nil {
			panic(err)
		}
		r = reflect.ValueOf(i)
	case t == "map":
		var m map[string]interface{}
		err := j.Decode(&m)
		if err != nil {
			panic(err)
		}
		r = reflect.ValueOf(m)
	case t == "string":
		var s string
		err := j.Decode(&s)
		if err != nil {
			panic(err)
		}
		r = reflect.ValueOf(s)
	case strings.HasPrefix(t, "st:"):
		var m map[string]interface{} // extract struct as map
		err := j.Decode(&m)
		if err != nil {
			panic(err)
		}
		r = reflect.ValueOf(m)
	case strings.HasPrefix(t, "sl:"):
		var a []interface{} // extract slice as array of interface{}
		err := j.Decode(&a)
		if err != nil {
			panic(err)
		}
		r = reflect.ValueOf(a)
	default:
		panic(fmt.Sprintf("cannot read/decode '%s' for caching", t))
	}

	return r
}
