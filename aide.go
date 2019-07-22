package fuel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var separator = ","

type Aide struct {
	Request  *http.Request
	Response http.ResponseWriter

	// variables extracted from http request
	query map[string]string
	post  map[string]string
	body  []byte
}

func (a *Aide) Query() map[string]string {

	// If not null => extraction has been done.
	// So simply return it.
	if a.query != nil {
		return a.query
	}

	a.query = make(map[string]string)
	qstring := a.Request.URL.Query()
	for k := range qstring {
		a.query[k] = strings.Join(qstring[k], separator)
	}
	return a.query
}

func (a *Aide) Post() map[string]string {

	// If not null => already parsed.
	// So just return
	if a.post != nil {
		return a.post
	}

	if a.Request.Method == http.MethodPost ||
		a.Request.Method == http.MethodPut ||
		a.Request.Method == http.MethodPatch {

		a.post = make(map[string]string)
		contentType := a.Request.Header.Get("Content-Type")
		switch {
		case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
			a.Request.ParseForm()
			for k := range a.Request.PostForm {
				a.post[k] = strings.Join(a.Request.PostForm[k], separator)
			}
		case strings.HasPrefix(contentType, "multipart/form-data"):
			// ParseMultiPart form should ideally populate
			// r.PostForm, but instead it fills r.Form
			// https://github.com/golang/go/issues/9305
			a.Request.ParseMultipartForm(1024 * 1024)
			for k := range a.Request.PostForm {
				a.post[k] = strings.Join(a.Request.PostForm[k], separator)
			}
		case strings.HasPrefix(contentType, "application/json"):
			if a.body == nil {
				a.body = getBody(a.Request)
			}
			var str = string(a.body)
			jsn := make(map[string]interface{})
			err := json.Unmarshal([]byte(str), &jsn)
			if err == nil {
				for key, data := range jsn {
					if val, ok := data.(string); ok { // put string value directly
						a.post[key] = val
					} else { // marshall non-string value before setting it
						byt, err := json.Marshal(data)
						if err != nil {
							panic(err)
						}
						a.post[key] = string(byt)
					}
				}
			} else {
				panic(err)
			}
		default:
			fmt.Println("Unhandled Content Type:", contentType)
			a.body = getBody(a.Request)
		}

		return a.post
	}

	return nil
}

func getBody(r *http.Request) []byte {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return b
}
