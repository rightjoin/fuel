package fuel

import (
	"fmt"
	"reflect"
	"strconv"
)

var faultSymbol = typeSymbol(reflect.TypeOf(Fault{}))

type Fault struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Inner   error  `json:"inner"`
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

func (f Fault) MarshalJSON() ([]byte, error) {

	// TODO: buffered string

	b := "{" +
		fmt.Sprintf(`"code": %d,`, f.Code) +
		fmt.Sprintf(` "message": %s `, strconv.Quote(f.Message))

	if f.Inner != nil {
		b += fmt.Sprintf(` ,"inner": %s`, strconv.Quote(f.Inner.Error()))
	}
	b += "}"

	return []byte(b), nil
}
