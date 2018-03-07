package fuel

import (
	"reflect"
)

type invoker struct {
	object     interface{}
	methodName string

	// keep it ready for dynamic invoking
	methodHandle reflect.Value
	inpSymbol    []string
	inpType      []reflect.Type
	outSymbol    []string
	outType      []reflect.Type
}

func newInvoker(obj interface{}, method string) invoker {

	ctype := reflect.TypeOf(obj)
	cvalue := reflect.ValueOf(obj)

	out := invoker{
		object:     obj,
		methodName: method,
	}

	found := false
	if ctype.Kind() == reflect.Ptr {
		// if the underlying method is defined
		// on address: func (s *struct) MethodName()
		if _, found = ctype.MethodByName(method); found {
			out.methodHandle = cvalue.MethodByName(method)
		}
	} else {
		// since the method was not found, it may have been defined
		// on the struct directly: func (s struct) MethodName()
		if _, found = ctype.MethodByName(method); found {
			out.methodHandle = cvalue.MethodByName(method)
		}
	}

	// if still not found, then there an issue
	if !found {
		panic("method not found: " + method + " in " + ctype.Name())
	}

	out.inpSymbol, out.inpType = func() ([]string, []reflect.Type) {
		count := out.methodHandle.Type().NumIn() // skip first parameter (me)
		inpS := make([]string, count)
		inpT := make([]reflect.Type, count)
		for i := 1; i < count; i++ {
			inpT[i-1] = out.methodHandle.Type().In(i)
			inpS[i-1] = typeSymbol(inpT[i-1])
		}
		return inpS, inpT
	}()

	out.outSymbol, out.outType = func() ([]string, []reflect.Type) {
		count := out.methodHandle.Type().NumOut()
		outS := make([]string, count)
		outT := make([]reflect.Type, count)
		for i := 0; i < count; i++ {
			outT[i] = out.methodHandle.Type().Out(i)
			outS[i] = typeSymbol(outT[i])
		}
		return outS, outT
	}()

	return out
}

func (i *invoker) invoke(params ...interface{}) []reflect.Value {
	inp := make([]reflect.Value, len(params))
	for i, p := range params {
		inp[i] = reflect.ValueOf(p)
	}
	return i.methodHandle.Call(inp)
}
