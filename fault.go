package fuel

import (
	"fmt"
	"reflect"
	"strconv"
)

var faultSymbol = typeSymbol(reflect.TypeOf(Fault{}))

type Fault struct {
	HTTPCode int    `json:"http_code"`
	ErrorNum int    `json:"error_num"`
	Message  string `json:"message"`
	Inner    error  `json:"inner"`
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
		// fmt.Sprintf(`"http_code": %d,`, f.HTTPCode) + No need to pass HTTP CODE to Clients
		fmt.Sprintf(`"error_num": %d,`, f.ErrorNum) +
		fmt.Sprintf(` "message": %s `, strconv.Quote(f.Message))

	if f.Inner != nil {
		b += fmt.Sprintf(` ,"inner": %s`, strconv.Quote(f.Inner.Error()))
	}
	b += "}"

	return []byte(b), nil
}
