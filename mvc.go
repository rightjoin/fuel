package fuel

import (
	"reflect"
)

type MvcOpts struct {
	Layout string
	Views  string
}

func defaultMvcOpts() MvcOpts {
	return MvcOpts{
		Views: "views",
	}
}

type View struct {
	// public members
	View   string
	Layout string
	Data   map[string]interface{}
}

var viewSymbol = typeSymbol(reflect.TypeOf(View{}))
