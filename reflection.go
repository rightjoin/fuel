package fuel

import (
	"fmt"
	"reflect"
)

func composedOf(item interface{}, parent interface{}) bool {

	it := reflect.TypeOf(item)
	if it.Kind() == reflect.Ptr {
		it = it.Elem()
	}
	if it.Kind() != reflect.Struct {
		panic("item must be struct")
	}

	pt := reflect.TypeOf(parent)
	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	if pt.Kind() != reflect.Struct {
		panic("parent must be struct")
	}

	// lookup field with parent's name
	f, ok := it.FieldByName(pt.Name())
	if !ok {
		return false
	}

	if !f.Anonymous {
		return false
	}

	if !f.Type.ConvertibleTo(pt) {
		return false
	}

	return true
}

func typeSymbol(t reflect.Type) (sym string) {

	switch t.Kind() {
	case reflect.Ptr:
		sym = "*" + typeSymbol(t.Elem())
	case reflect.Map:
		sym = "map"
	case reflect.Struct:
		sym = fmt.Sprintf("st:%s.%s", t.PkgPath(), t.Name())
	case reflect.Interface:
		sym = fmt.Sprintf("i:%s.%s", t.PkgPath(), t.Name())
	case reflect.Array, reflect.Slice:
		sym = fmt.Sprintf("sl:%s.%s", t.PkgPath(), t.Name())
	default:
		sym = t.Name()
	}

	return
}
