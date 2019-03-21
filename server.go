package fuel

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/unrolled/render"

	"github.com/gorilla/mux"
	"github.com/rightjoin/rutl/conv"
	"github.com/rightjoin/stak"
)

const defaultPort = 8080

type Server struct {
	http.Server
	Fixture
	Port      int
	mux       *mux.Router
	svcs      []serviceComposite
	endpoints map[string]endpoint
	middle    map[string]func(http.Handler) http.Handler
	caches    map[string]stak.Cache

	_MvcOptions MvcOpts // hidden for now (NO MVC)
}

func NewServer() Server {
	return Server{
		Server: http.Server{
			ReadTimeout:  1 * time.Minute,
			WriteTimeout: 1 * time.Minute,
			IdleTimeout:  15 * time.Minute,
		},

		Port:        defaultPort,
		mux:         mux.NewRouter(),
		middle:      make(map[string]func(http.Handler) http.Handler),
		svcs:        make([]serviceComposite, 0),
		endpoints:   make(map[string]endpoint),
		caches:      make(map[string]stak.Cache, 0),
		_MvcOptions: defaultMvcOpts(),
	}
}

func (s *Server) DefineMiddleware(name string, fn func(http.Handler) http.Handler) {
	if _, ok := s.middle[name]; ok {
		panic("middleware already defined: " + name)
	}
	s.middle[name] = fn
}

func (s *Server) DefineCache(name string, c stak.Cache) {
	if _, ok := s.caches[name]; ok {
		panic("cache already defined: " + name)
	}
	s.caches[name] = c
}

func (s *Server) AddService(svc serviceComposite) {

	// input validation:
	if reflect.TypeOf(svc).Kind() != reflect.Ptr {
		panic("service must be passed as a pointer")
	}
	if !composedOf(svc, Service{}) {
		panic("service must be composed of fuel.Service")
	}

	// store it
	s.svcs = append(s.svcs, svc)
}

func (s *Server) loadEndpoints() {

	for _, svc := range s.svcs {

		// load endpoints::

		ctype := reflect.TypeOf(svc).Elem()
		cvalue := reflect.ValueOf(svc).Elem()

		// build service fixture
		fixContr := func() Fixture {
			fldCont, _ := ctype.FieldByName("Service")
			fixTag := newFixture(fldCont.Tag)
			fixCode := cvalue.FieldByName("Service").FieldByName("Fixture").Interface().(Fixture)
			fixCode.Parent = &fixTag
			fixTag.Parent = &s.Fixture

			// if there is no root value set, then
			// use service name (minus -service) as Root
			if fixCode.getRoot() == "" {
				fixCode.Root = conv.CaseURL(ctype.Name())
				if strings.HasSuffix(fixCode.Root, "-service") {
					fixCode.Root = fixCode.Root[0 : len(fixCode.Root)-len("-service")]
				}
			}
			return fixCode
		}()

		for i := 0; i < ctype.NumField(); i++ {
			fieldType := ctype.FieldByIndex([]int{i})
			fieldValue := cvalue.FieldByIndex([]int{i})

			//fmt.Println(">>>>>>", typeSymbol(fieldType.Type), fieldType.Type.String())

			// must be of given http return methods, else skip
			//method := fieldType.Type.String()[len("fuel")+1:]
			switch fieldType.Type.String() {
			case "fuel.GET", "fuel.PUT", "fuel.POST", "fuel.DELETE":
			default:
				continue
			}

			// build endpoint fixture
			fix := func() Fixture {
				var out Fixture
				fixTag := newFixture(fieldType.Tag)
				fixTag.Parent = &fixContr
				out = fixTag

				if fieldType.Name[0:1] == strings.ToUpper(fieldType.Name[0:1]) {
					// this field starts Uppercase,
					// so extract GET.Fixture also
					// (when lowercase this is not possible as the field value
					//  become unexported)
					fixCode := fieldValue.FieldByName("Fixture").Interface().(Fixture)
					fixCode.Parent = &fixTag
					out = fixCode
				}

				// if no route defined, then assume it basis
				// the name of the action
				if out.getRoute() == "" {
					out.Route = conv.CaseURL(fieldType.Name)
				}

				return out
			}()

			// build the endpoint, and store it in the server
			epoint := newEndpoint(fix, svc, fieldType, s)
			uniqURL := epoint.uniqueURL()
			if _, ok := s.endpoints[uniqURL]; ok {
				panic("cannot use same url again: " + uniqURL)
			}
			s.endpoints[uniqURL] = epoint

			// print it in formatted manner
			spaces := "  "
			for i := strings.Index(uniqURL, ":"); i < len("DELETE:")-1; i++ {
				spaces += " "
			}
			fmt.Println(spaces + uniqURL)
		}
	}
}

func (s *Server) Run() {

	// load endpoints:
	// we do it at the end in the 'Run' step because
	// the user may add a cache later on (after calling AddService),
	// or the user may add middleware later on
	s.loadEndpoints()

	// setup the renderer:
	// basis some of the settings passed to server's MvcOptions
	rndr = render.New(render.Options{
		Directory:  s._MvcOptions.Views,
		Layout:     s._MvcOptions.Layout,
		Extensions: []string{".html"},
	})

	s.Addr = fmt.Sprintf(":%d", s.Port)
	s.Server.Handler = s.mux
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

var CacheKey = func(r *http.Request) string {
	return r.RequestURI
}
