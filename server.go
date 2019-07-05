package fuel

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/logrusorgru/aurora"
	"github.com/rightjoin/rutl/conv"
	"github.com/rightjoin/stak"
	"github.com/unrolled/render"
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

		// Build service fixture as follows:
		//  1. Higest Parent = Server
		//  2. -1 Parent = Service configuration programmatically
		//  3. -2 Parent = Service configuration declaratively
		// Note: We give higher precedence to declarative configuration
		// over programmatic configuration
		fixParent := func() Fixture {
			field, _ := ctype.FieldByName("Service")
			fixTag := newFixture(field.Tag)
			fixCode := cvalue.FieldByName("Service").FieldByName("Fixture").Interface().(Fixture)

			// Set parent hierarchy
			fixTag.Parent = &fixCode
			fixCode.Parent = &s.Fixture

			// If there is no root value set, then
			// use service name (minus -service) as Root
			if fixTag.getRoot() == "" {
				extractRoot := conv.CaseURL(ctype.Name())
				if strings.HasSuffix(extractRoot, "-service") {
					extractRoot = extractRoot[0 : len(extractRoot)-len("-service")]
				}
				fixTag.Root = extractRoot
			}
			return fixTag
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

			// Build endpoint fixture.
			// Again, we give higher precedence to
			// declarative configuration over programmatic configuration
			fix := func() Fixture {
				fixTag := newFixture(fieldType.Tag)
				fixTag.Parent = &fixParent

				// For uppercase fields, exract programmatic configuration also
				// (when lowercase this is not possible as the field value
				//  become unexported)
				if fieldType.Name[0:1] == strings.ToUpper(fieldType.Name[0:1]) {
					// extract GET.Fixture /etc
					fixCode := fieldValue.FieldByName("Fixture").Interface().(Fixture)
					fixCode.Parent = &fixParent
					fixTag.Parent = &fixCode
				}

				// if no route defined, then assume it basis
				// the name of the action
				if fixTag.Route == "" {
					fixTag.Route = conv.CaseURL(fieldType.Name)
				}

				return fixTag
			}()

			// Build the endpoint, and store it in the server
			epoint := newEndpoint(fix, svc, fieldType, s)
			uniqURL := epoint.uniqueURL()
			if _, ok := s.endpoints[uniqURL]; ok {
				panic("cannot use same url again: " + uniqURL)
			}
			s.endpoints[uniqURL] = epoint

			// Print it in formatted manner
			spaces := "  "
			for i := strings.Index(uniqURL, ":"); i < len("DELETE:")-1; i++ {
				spaces += " "
			}
			fmt.Println(spaces, aurora.Blue(epoint.method()), epoint.Fixture.getURL())
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

	// Print port information
	fmt.Println("    ", aurora.Black("Port").BgGreen(), aurora.Green(s.Port))

	s.Addr = fmt.Sprintf(":%d", s.Port)
	s.Server.Handler = s.mux
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

var mx = sync.Mutex{}
var portToUse = 9595

// RunTestInstance selects a port for the underlying server to be run.
// It starts at 9595, and with each invocation picks the next one.
// It also starts the server async and returns the root localhost
// url along with the port number
func (s *Server) RunTestInstance() (url string, port int) {
	// Critical Section
	mx.Lock()
	{
		portToUse++
		s.Port = portToUse
	}
	mx.Unlock()

	go s.Run()

	// Give the server 50ms to fire up
	time.Sleep(50 * time.Millisecond)

	return fmt.Sprintf("http://localhost:%d", s.Port), s.Port
}

var CacheKey = func(r *http.Request) string {
	return r.RequestURI
}
