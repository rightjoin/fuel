package fuel

import "net/http"

type Aide struct {
	Request  *http.Request
	Response http.ResponseWriter

	// variables extracted from http request

}
