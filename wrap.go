package fuel

import "fmt"

type BodyWrap interface {
	SetData(interface{})
	SetError(error)
	SetFault(Fault)
}

type ApiResponse struct {
	Data     interface{}  `json:"data"`
	HasError bool         `json:"has_error"`
	Errors   []CodedError `json:"errors"`
}

type CodedError struct {
	Code         int
	ErrorMessage string
}

func (c CodedError) Error() string {
	return c.ErrorMessage
}

func (api *ApiResponse) SetData(data interface{}) {
	api.Data = data
}

func (api *ApiResponse) SetError(e error) {
	api.HasError = true
	if api.Errors == nil {
		api.Errors = []CodedError{}
	}
	api.Errors = append(api.Errors, CodedError{
		ErrorMessage: e.Error(),
	})

	fmt.Println("setting errors", api.Errors)
}

func (api *ApiResponse) SetFault(f Fault) {
	api.HasError = true
	if api.Errors == nil {
		api.Errors = []CodedError{}
	}
	api.Errors = append(api.Errors, CodedError{
		Code:         f.ErrorNum,
		ErrorMessage: f.Message,
	})

	fmt.Println("setting faults", api.Errors)
}
