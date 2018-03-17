package fuel

import (
	"reflect"
)

type invoker struct {
	object     interface{}
	methodName string

	// keep it ready for dynamic invoking
	methodType   reflect.Method
	methodHandle reflect.Value
	inpSymbol    []string
	inpType      []reflect.Type
	outSymbol    []string
	outType      []reflect.Type
}

func newInvoker(obj interface{}, method string) invoker {

	ctype := reflect.TypeOf(obj)
	cvalue := reflect.ValueOf(obj)

	invk := invoker{
		object:     obj,
		methodName: method,
	}

	found := false
	if ctype.Kind() == reflect.Ptr {
		// if the underlying method is defined
		// on address: func (s *struct) MethodName()
		if invk.methodType, found = ctype.MethodByName(method); found {
			invk.methodHandle = cvalue.MethodByName(method)
		}
	} else {
		// since the method was not found, it may have been defined
		// on the struct directly: func (s struct) MethodName()
		if invk.methodType, found = ctype.MethodByName(method); found {
			invk.methodHandle = cvalue.MethodByName(method)
		}
	}

	// if still not found, then there is an issue
	if !found {
		panic("method not found: " + method + " in " + ctype.Name())
	}

	invk.inpSymbol, invk.inpType = func() ([]string, []reflect.Type) {
		count := invk.methodType.Type.NumIn() - 1 // skip first parameter (me)
		inpS := make([]string, count)
		inpT := make([]reflect.Type, count)
		for i := 1; i < count+1; i++ {
			inpT[i-1] = invk.methodType.Type.In(i)
			inpS[i-1] = typeSymbol(inpT[i-1])
		}
		return inpS, inpT
	}()

	invk.outSymbol, invk.outType = func() ([]string, []reflect.Type) {
		count := invk.methodType.Type.NumOut()
		outS := make([]string, count)
		outT := make([]reflect.Type, count)
		for i := 0; i < count; i++ {
			outT[i] = invk.methodHandle.Type().Out(i)
			outS[i] = typeSymbol(outT[i])
		}
		return outS, outT
	}()

	return invk
}

func (i *invoker) invoke(params ...interface{}) []reflect.Value {
	inp := make([]reflect.Value, len(params))
	for i, p := range params {
		inp[i] = reflect.ValueOf(p)
	}
	return i.methodHandle.Call(inp)
}
