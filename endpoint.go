package fuel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

type endpoint struct {
	Fixture
	controller service
	field      reflect.StructField
	invoker

	standardHandler bool
	usesAide        bool
}

func newEndpoint(fix Fixture, contr service, fld reflect.StructField, server *Server) endpoint {

	// if field is abc, then methods must be Abc
	// if field is Abc, then method must be Abc_
	var seekMethod string
	if fld.Name[0:1] != strings.ToUpper(fld.Name[0:1]) {
		seekMethod = strings.ToUpper(fld.Name[0:1]) + fld.Name[1:]
	} else {
		seekMethod = fld.Name + "_"
	}
	var inv = newInvoker(contr, seekMethod)

	var aide = typeSymbol(reflect.TypeOf(Aide{}))
	out := endpoint{
		Fixture:    fix,
		controller: contr,
		field:      fld,
		invoker:    inv,

		standardHandler: func() bool {
			return len(inv.inpSymbol) == 2 &&
				inv.inpSymbol[0] == "i:net/http.ResponseWriter" &&
				inv.inpSymbol[1] == "*st:net/http.Request" &&
				len(inv.outSymbol) == 0
		}(),
		usesAide: func() bool {
			for i := 0; i < len(inv.inpSymbol)-1; i++ { // not present at any other index
				if inv.inpSymbol[i] == aide {
					panic("Aide parameter must come at end: " + seekMethod)
				}
			}
			return len(inv.inpSymbol) > 0 && inv.inpSymbol[len(inv.inpSymbol)-1] == aide // present at last index
		}(),
	}

	// validations::

	// length of mux variables should EQUAL input count of function
	muxVars := extractMuxVars(out.getURL())
	if !out.standardHandler {
		count := len(out.inpSymbol)
		if out.usesAide {
			count--
		}
		if len(muxVars) != count {
			title := fmt.Sprintf("%d inputs of url (%s) do not match %d of func:%s", len(muxVars), out.getURL(), count, seekMethod)
			panic(title)
		}
	}

	// function inputs should be supported type
	var supp = []string{"int", "uint", "string", aide}
	if !out.standardHandler {
		for _, inp := range out.inpSymbol {
			match := false
			for _, sup := range supp {
				if sup == inp {
					match = true
					break
				}
			}
			if !match {
				panic("func input params must be: " + strings.Join(supp[0:len(supp)-1], "|"))
			}
		}
	}

	// function outputs must of right format
	if !out.standardHandler {
		switch len(inv.outSymbol) {
		case 1:
			if !acceptableOutput(inv.outSymbol[0]) {
				panic("incorrect/unsupported output param in: " + seekMethod)
			}
		case 2:
			if !acceptableOutput(inv.outSymbol[0]) {
				panic("incorrect/unsupported output param in: " + seekMethod)
			}
			if inv.outSymbol[1] != "i:.error" {
				panic("second output param must be error: " + seekMethod)
			}
		default:
			panic("cannot have more than two return params: " + seekMethod)
		}
	}

	// finally, setup mux handlers
	out.setupMuxHandlers(server)

	return out
}

func (e *endpoint) setupMuxHandlers(server *Server) {

	m := interpose.New()
	if e.getMiddleware() != nil {
		for _, midw := range e.getMiddleware() {
			m.Use(server.middle[midw])
		}
	}
	fn := processRequest(e)
	m.UseHandler(http.HandlerFunc(fn))

	server.mux.Handle(e.getURL(), m).Methods(e.method())
}

func (e *endpoint) method() string {
	return e.field.Type.String()[len("fuel")+1:]
}

func (e *endpoint) uniqueURL() string {
	return e.method() + ":" + e.Fixture.getURL()
}

func processRequest(e *endpoint) func(http.ResponseWriter, *http.Request) {

	// call predefined function if the
	// handler is a standard one
	if e.standardHandler {
		return func(w http.ResponseWriter, r *http.Request) {
			e.invoke(w, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// get the inputs that need to be passed to the underlying handler
		inputs := func() []interface{} {
			muxVals := mux.Vars(r)
			inp := make([]interface{}, len(muxVals))
			i := 0
			var err error
			for _, val := range muxVals {
				switch e.inpSymbol[i] {
				case "string":
					inp[i] = val
				case "int":
					inp[i], err = strconv.Atoi(val)
					if err != nil {
						panic("input param expects 'int': " + val)
					}
				case "uint":
					inp[i], err = strconv.ParseUint(val, 10, 32)
					if err != nil {
						panic("input param expects 'uint': " + val)
					}
				}
				i++
			}
			return inp
		}()

		// do we need aide?
		if e.usesAide {
			inputs = append(inputs, Aide{Request: r, Response: w})
		}

		outputs := e.invoke(inputs...)
		writeHTTP(e, w, r, outputs)
	}
}

func writeHTTP(e *endpoint, w http.ResponseWriter, r *http.Request, data []reflect.Value) {
	if len(data) == 1 {
		writeItem(e, w, r, data[0])
	} else {

	}
}

func writeItem(e *endpoint, w http.ResponseWriter, r *http.Request, item reflect.Value) {

	// if reflect value is a ptr, then lets
	// just process its internal element
	runtimeType := reflect.TypeOf(item.Interface())
	if runtimeType.Kind() == reflect.Ptr {
		writeItem(e, w, r, item.Elem())
		return
	}

	var symbol = typeSymbol(runtimeType)

	var sendJSON = func() {
		jsn, err := json.Marshal(item.Interface())
		if err != nil {

		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(jsn)))
		w.Write(jsn)
	}

	switch {
	case symbol == "string":
		{
			data := item.Interface().(string)
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Lenght", strconv.Itoa(len(data)))
			fmt.Fprintf(w, "%s", data)
		}
	case symbol == "map":
		// TODO: string -> interface{}
		sendJSON()
	case strings.HasPrefix(symbol, "st:"):
		sendJSON()
	case strings.HasPrefix(symbol, "sl:"):
		sendJSON()
	case symbol == "i.error":
	default:
		panic("unable to process: " + symbol)
	}
}
