package fuel

import (
	"reflect"
)

var faultSymbol = typeSymbol(reflect.TypeOf(Fault{}))

type Fault struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Inner   error  `json:"inner,omitempty"`
}

func (f Fault) Error() string {
	if f.Message != "" {
		return f.Message
	}

	if f.Inner != nil {
		return f.Inner.Error()
	}

	panic("Emtpy fault!")
}
