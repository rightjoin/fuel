package fuel

import (
	"net/http"
	"strings"
)

type Aide struct {
	Request  *http.Request
	Response http.ResponseWriter

	// variables extracted from http request
	query map[string]string
}

func (a *Aide) Query() map[string]string {

	// If not null => extraction has been done.
	// So simply return it
	if a.query != nil {
		return a.query
	}

	// For GET requests, parse Request.Form
	if a.Request.Method == http.MethodGet {
		a.query = make(map[string]string)
		for k := range a.Request.Form {
			a.query[k] = strings.Join(a.Request.Form[k], ",")
		}
		return a.query
	}

	return nil
}
