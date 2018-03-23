package fuel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rightjoin/stag"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

var cacheWriter = serializer{}

type endpoint struct {
	Fixture
	controller service
	field      reflect.StructField
	invoker
	paramName []string

	standardHandler bool
	usesAide        bool

	myCache    stag.Cache
	myCacheDur time.Duration
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

	var inv = invoker{}
	if fix.getStub() == "" {
		// we setup an actual invoker only if it is not a stub.
		// otherwise newInvoker panics as it does not find the
		// implementation method.
		inv = newInvoker(contr, seekMethod)
	}

	var aide = typeSymbol(reflect.TypeOf(Aide{}))
	out := endpoint{
		Fixture:    fix,
		controller: contr,
		field:      fld,
		invoker:    inv,
		paramName: func() []string {
			output := make([]string, 0)
			dirs := strings.Split(fix.getURL(), "/")
			for _, dir := range dirs {
				if strings.HasPrefix(dir, "{") && strings.HasSuffix(dir, "}") {
					p := dir[1 : len(dir)-1]
					if strings.Contains(p, ":") {
						p = p[0:strings.Index(p, ":")]
					}
					output = append(output, p)
				}
			}
			return output
		}(),

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

	// caching ::
	out.myCache, out.myCacheDur = func() (stag.Cache, time.Duration) {
		name := fix.getCache()
		ttl := fix.getTTL()
		dur, err := time.ParseDuration(ttl)
		if ttl != "" && err != nil {
			panic("incorrect ttl: " + ttl)
		}
		if name != "" {
			c, found := server.caches[name]
			if found {
				return c, dur
			}
			panic("cache provider not found: " + name)
		}
		return nil, 0
	}()

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
	if !out.standardHandler && fix.getStub() == "" {
		switch len(inv.outSymbol) {
		case 1:
			if !acceptableOutput(inv.outSymbol[0]) {
				panic("incorrect or unsupported output param in: " + seekMethod)
			}
		case 2:
			if !acceptableOutput(inv.outSymbol[0]) {
				panic("incorrect or unsupported output param in: " + seekMethod)
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
		for _, midName := range e.getMiddleware() {
			midw, found := server.middle[midName]
			if !found {
				panic("middleware not found: " + midName)
			}
			m.Use(midw)
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

	// stub::
	// parse file contents and serve it back
	if e.Stub != "" {
		return func(w http.ResponseWriter, r *http.Request) {
			data, err := readFile(e.Stub)
			if err == nil {
				fmt.Fprintf(w, "%s", data)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Eror reading stub: %s", e.Stub)
		}
	}

	// standard handler::
	// call predefined function if it is a standard handler
	if e.standardHandler {
		return func(w http.ResponseWriter, r *http.Request) {
			e.invoke(w, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// get the inputs that need to be passed to the underlying handler
		params := make([]interface{}, 0)
		if len(e.paramName) > 0 {
			var err error
			params = make([]interface{}, len(e.paramName))
			muxVals := mux.Vars(r)
			for i, pName := range e.paramName {
				val := muxVals[pName]
				switch e.inpSymbol[i] {
				case "string":
					params[i] = val
				case "int":
					params[i], err = strconv.Atoi(val)
					if err != nil {
						panic("input param expects 'int': " + val)
					}
				case "uint":
					u, err := strconv.ParseUint(val, 10, 32)
					if err != nil {
						panic("input param expects 'uint': " + val)
					}
					params[i] = uint(u)
					// TODO: default?
				}
			}

		}

		// do we need aide?
		if e.usesAide {
			params = append(params, Aide{Request: r, Response: w})
		}

		var cacheOn = e.myCacheDur > 0 && r.Method == http.MethodGet
		var outputs []reflect.Value

		// cached vs non-cached behavior
		if !cacheOn {
			// there is no caching
			outputs = e.invoke(params...)
		} else {
			// try finding cached value
			val, err := e.myCache.Get(CacheKey(r))
			if err == nil {
				// hit hit hit!
				outputs = cacheWriter.read(bytes.NewBuffer(val), e.outType)
			} else {
				// invoke the normal method
				outputs = e.invoke(params...)
				// try saving to cache
				buf := cacheWriter.write(outputs)
				e.myCache.Set(CacheKey(r), buf.Bytes(), e.myCacheDur)
				// TODO: better error handling:
				// - don't write to cache again for 5 sec
				// - don't read from cache until that time etc
			}
		}

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
