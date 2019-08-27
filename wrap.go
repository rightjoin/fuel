package fuel

type BodyWrap interface {
	SetData(interface{})
	SetError(error)
	SetFault(Fault)
}

type ApiResponse struct {
	Data    interface{}  `json:"data"`
	Success bool         `json:"success"`
	Errors  []CodedError `json:"errors"`
}

type CodedError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_message"`
}

func (c CodedError) Error() string {
	return c.ErrorMessage
}

func (api *ApiResponse) SetData(data interface{}) {
	api.Data = data
}

func (api *ApiResponse) SetError(e error) {
	api.Success = false
	if api.Errors == nil {
		api.Errors = []CodedError{}
	}
	api.Errors = append(api.Errors, CodedError{
		ErrorMessage: e.Error(),
	})
}

func (api *ApiResponse) SetFault(f Fault) {
	api.Success = false
	if api.Errors == nil {
		api.Errors = []CodedError{}
	}
	api.Errors = append(api.Errors, CodedError{
		Code:         f.ErrorNum,
		ErrorMessage: f.Message,
	})
}
