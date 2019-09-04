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

	"github.com/asaskevich/govalidator"
	newrelic "github.com/newrelic/go-agent"

	"github.com/rightjoin/rutl/conv"
	"github.com/rightjoin/stak"
	"github.com/unrolled/render"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

var cacheWriter = serializer{}

type endpoint struct {
	Fixture
	parent serviceComposite
	field  reflect.StructField
	invoker
	paramName []string

	standardHandler bool
	usesAide        bool
	fillStruct      bool

	myCache    stak.Cache
	myCacheDur time.Duration

	wrapperFn func() BodyWrap

	mvcOptions MvcOpts
	viewDir    string
}

func newEndpoint(fix Fixture, myparent serviceComposite, fld reflect.StructField, server *Server) endpoint {

	// if field is abc, then methods must be Abc
	// if field is Abc, then method must be Abc_
	var seekMethod string
	if fld.Name[0:1] != strings.ToUpper(fld.Name[0:1]) {
		seekMethod = strings.ToUpper(fld.Name[0:1]) + fld.Name[1:]
	} else {
		seekMethod = fld.Name + "_"
	}

	var inv = invoker{}
	var skipInvoker = fix.getStub() != "" || fix.getFolder() != ""
	if !skipInvoker {
		// we setup an actual invoker only if it is not a stub,
		// and also it is not a static file server.
		// otherwise newInvoker panics as it does not find the
		// implementation method.
		inv = newInvoker(myparent, seekMethod)
	}

	var aide = typeSymbol(reflect.TypeOf(Aide{}))
	out := endpoint{
		Fixture: fix,
		parent:  myparent,
		field:   fld,
		invoker: inv,
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
		fillStruct: func() bool {
			ln := len(inv.inpSymbol)
			if ln == 0 {
				return false
			}
			// If last param is aide, then second last must be checked
			// Otherwise last should be checked
			hasAide := ln > 0 && inv.inpSymbol[ln-1] == aide
			meth := fld.Type.String()[len("fuel")+1:]
			loop := ln
			if hasAide {
				loop = ln - 1
			}
			for i := 0; i < loop; i++ {
				if meth != "POST" && meth != "PUT" && strings.HasPrefix(inv.inpSymbol[i], "st:") {
					panic("Struct can only be used with PUT and POST : " + seekMethod)
				}
			}

			if hasAide {
				if ln == 1 {
					return false
				}
				return strings.HasPrefix(inv.inpSymbol[ln-2], "st:")
			}
			return strings.HasPrefix(inv.inpSymbol[ln-1], "st:")
		}(),
		mvcOptions: server._MvcOptions,
		viewDir: func() string {
			name := reflect.TypeOf(myparent).Elem().Name()
			snake := conv.CaseURL(name)
			if strings.HasSuffix(snake, "-service") {
				return snake[0 : len(name)-len("-service")+1]
			}
			return snake
		}(),
		wrapperFn: server.ResponseFormat,
	}

	// caching ::
	out.myCache, out.myCacheDur = func() (stak.Cache, time.Duration) {
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
		if out.fillStruct {
			count--
		}
		if len(muxVars) != count {
			title := fmt.Sprintf("%d inputs of url (%s) do not match %d of func:%s", len(muxVars), out.getURL(), count, seekMethod)
			panic(title)
		}
	}

	// function inputs should be supported type
	var supp = []string{"int", "uint", "string"}
	if !out.standardHandler {
		ln := len(out.inpSymbol) - 1
		if out.usesAide {
			ln--
		}
		if out.fillStruct {
			ln--
		}
		for i := 0; i <= ln; i++ {
			inp := out.inpSymbol[i]
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
	var dontValidate = out.standardHandler || fix.getStub() != "" || fix.getFolder() != ""
	if !dontValidate {
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
			panic("cannot have zero or more than two return params: " + seekMethod)
		}
	}

	// finally, setup mux handlers
	out.setupMuxHandlers(server)

	return out
}

func (e *endpoint) setupMuxHandlers(server *Server) {

	m := interpose.New()
	mw := e.getMiddleware()
	skipMware := mw == nil || (len(mw) == 1 && mw[0] == "-")
	if !skipMware {
		for _, midName := range mw {
			midw, found := server.middle[midName]
			if !found {
				panic("middleware not found: " + midName)
			}
			m.Use(midw)
		}
	}

	if e.getFolder() != "" {
		// setup static server
		path := e.getURL()
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
		//muxer.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./files/"))))
		m.UseHandler(http.StripPrefix(path, http.FileServer(http.Dir(e.getFolder()))))
		server.mux.PathPrefix(path).Handler(m).Methods(e.method())
	} else {
		// normal request processing is handled through
		// processRequest function
		fn := processRequest(e)
		m.UseHandler(http.HandlerFunc(fn))

		if server.NewRelicApp != nil {
			_, h := newrelic.WrapHandle(*server.NewRelicApp, e.getURL(), m.Handler())
			server.mux.Handle(e.getURL(), h).Methods(e.method())
		} else {
			server.mux.Handle(e.getURL(), m).Methods(e.method())
		}
	}

}

func (e *endpoint) method() string {
	return e.field.Type.String()[len("fuel")+1:]
}

func (e *endpoint) uniqueURL() string {
	return e.method() + ":" + e.Fixture.getURL()
}

func processRequest(e *endpoint) func(http.ResponseWriter, *http.Request) {

	// stub :
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

		// Do we need a wrapper?
		var wrap BodyWrap
		if e.Fixture.getWrap() == "true" && e.wrapperFn != nil {
			if len(r.Header.Get("No-Wrap")) == 0 {
				wrap = e.wrapperFn()
			}
		}

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
		var ad = &Aide{Request: r, Response: w}
		if e.fillStruct {
			ln := len(e.inpType)
			var t reflect.Type
			if e.usesAide {
				t = e.inpType[ln-2]
			} else {
				t = e.inpType[ln-1]
			}
			// fmt.Println(e.inpSymbol)
			// fmt.Println(refl.Signature(t))
			var p = reflect.New(t)
			var addr = p.Interface()
			buf, err := json.Marshal(ad.Post())
			if err != nil {
				writeItem(e, w, r, reflect.ValueOf(err), wrap)
				return
			}
			// fmt.Println("Marshalled To:", string(buf))
			err = json.Unmarshal(buf, addr)
			if err != nil {
				writeItem(e, w, r, reflect.ValueOf(err), wrap)
				return
			}
			ok, err := govalidator.ValidateStruct(addr)
			// fmt.Println(ok, err)
			if !ok {
				writeItem(e, w, r, reflect.ValueOf(err), wrap)
				return
			}
			params = append(params, p.Elem().Interface())
		}

		if e.usesAide {
			params = append(params, *ad)
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

		writeHTTP(e, w, r, outputs, wrap)
	}
}

func writeHTTP(e *endpoint, w http.ResponseWriter, r *http.Request, data []reflect.Value, wrap BodyWrap) {

	if len(data) == 1 {
		writeItem(e, w, r, data[0], wrap)
	} else {
		// Second parameter is type error.
		// If it is NIL then all good, so
		// process everything as normal.
		// Otherwise process error.
		if data[1].IsNil() {
			writeItem(e, w, r, data[0], wrap)
		} else {

			// Set content
			if !data[0].IsNil() {
				if wrap != nil {
					wrap.SetData(data[0].Interface())
				}
			}

			writeItem(e, w, r, data[1], wrap)
		}
	}
}

func writeItem(e *endpoint, w http.ResponseWriter, r *http.Request, item reflect.Value, wrap BodyWrap) {

	_, isError := item.Interface().(error)

	// If reflect value is a ptr, then
	// remove indirection (unless an error).
	// If we remove indirectin for error values too,
	// then it messes up their conversion to error.
	//runtimeType := reflect.TypeOf(item.Interface())
	runtimeType := item.Type()
	if !isError && runtimeType.Kind() == reflect.Ptr {
		fmt.Println("remove-indirection", typeSymbol(runtimeType))
		writeItem(e, w, r, item.Elem(), wrap)
		return
	}

	var symbol = typeSymbol(runtimeType)

	var sendJSON = func(status ...int) {

		// If no status passed to func
		// then use OK. Else used first value
		sendStatus := http.StatusOK
		if len(status) != 0 {
			sendStatus = status[0]
		}

		// TODO: error validation
		rndr.JSON(w, sendStatus, item.Interface())
	}

	var sendJSONOf = func(data interface{}, status ...int) {

		// If no status passed to func
		// then use OK. Else used first value
		sendStatus := http.StatusOK
		if len(status) != 0 && status[0] != 0 {
			sendStatus = status[0]
		}

		// TODO: error validation
		rndr.JSON(w, sendStatus, data)
	}

	//helper function for view rendering
	var renderView = func() {
		v := item.Interface().(View)

		// view
		var view = v.View
		if view == "" {
			view = e.field.Name
		}
		var finalView = cleanMultSlash(e.viewDir + "/" + view)

		// layout
		var layout = v.Layout
		if layout == "" {
			layout = e.mvcOptions.Layout
		}

		//fmt.Println(">>", e.mvcOptions.Views+"/"+e.viewDir+"/"+v.View, ">>", e.mvcOptions.Views+"/"+v.Layout)

		rndr.HTML(w, http.StatusOK, finalView, v.Data, render.HTMLOptions{
			Layout: layout,
		})
	}

	fmt.Println("writeItem() symbol->>", symbol)
	switch {
	case symbol == faultSymbol:
		f := item.Interface().(Fault)
		httpStatus := http.StatusExpectationFailed
		if f.HTTPCode != 0 {
			httpStatus = f.HTTPCode
		}
		if wrap == nil {
			sendJSON(httpStatus)
		} else {
			wrap.SetFault(f)
			sendJSONOf(wrap, httpStatus)
		}
	case isError || (symbol == "i:.error"):
		if item.Interface() == nil {
			if wrap == nil {
				success := map[string]interface{}{"success": 1}
				writeItem(e, w, r, reflect.ValueOf(success), wrap)
			} else {
				sendJSONOf(wrap)
			}
			return
		}
		f, faulty := item.Interface().(Fault)
		if faulty && wrap != nil {
			wrap.SetFault(f)
		}
		if !faulty {
			if wrap != nil {
				wrap.SetError(item.Interface().(error))
			}
			f = Fault{Message: "An error occurred", Inner: item.Interface().(error), ErrorNum: 9999}
			f.HTTPCode = http.StatusExpectationFailed
			fmt.Println("wrapping error into fault:", f.Inner, "; and outer =>", f)
		}
		if wrap == nil {
			writeItem(e, w, r, reflect.ValueOf(f), wrap)
		} else {
			if f.HTTPCode == 0 {
				f.HTTPCode = http.StatusExpectationFailed
			}
			sendJSONOf(wrap, f.HTTPCode)
		}
	case symbol == "string":
		{
			if wrap == nil {
				rndr.Text(w, http.StatusOK, item.Interface().(string))
			} else {
				wrap.SetData(item.Interface().(string))
				sendJSONOf(wrap)
			}

			// data := item.Interface().(string)
			// w.Header().Set("Content-Type", "text/plain")
			// w.Header().Set("Content-Lenght", strconv.Itoa(len(data)))
			// fmt.Fprintf(w, "%s", data)
		}
	case symbol == "map":
		// TODO: string -> interface{}
		if wrap == nil {
			sendJSON()
		} else {
			wrap.SetData(item.Interface())
			sendJSONOf(wrap)
		}
	case symbol == viewSymbol:
		renderView()
	case strings.HasPrefix(symbol, "st:"):
		if wrap == nil {
			sendJSON()
		} else {
			wrap.SetData(item.Interface())
			sendJSONOf(wrap)
		}
	case strings.HasPrefix(symbol, "sl:"):
		if wrap == nil {
			sendJSON()
		} else {
			wrap.SetData(item.Interface())
			sendJSONOf(wrap)
		}
	case symbol == "i:.":
		writeItem(e, w, r, reflect.ValueOf(item.Interface()), wrap)
	default:
		panic("unable to process: " + symbol)
	}
}
