package fuel

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rightjoin/txt"
)

const defaultPort = 8421

type Server struct {
	http.Server
	Fixture
	Port        int
	mux         *mux.Router
	middle      map[string]func(http.Handler) http.Handler
	controllers []service
	endpoints   map[string]endpoint
}

func NewServer() Server {
	return Server{
		Port:        defaultPort,
		mux:         mux.NewRouter(),
		middle:      make(map[string]func(http.Handler) http.Handler),
		controllers: make([]service, 0),
		endpoints:   make(map[string]endpoint),
	}
}

func (s *Server) DefineMiddleware(name string, fn func(http.Handler) http.Handler) {
	if _, ok := s.middle[name]; ok {
		panic("middleware already added: " + name)
	}
	s.middle[name] = fn
}

func (s *Server) AddController(controller service) {

	// input validation:
	if reflect.TypeOf(controller).Kind() != reflect.Ptr {
		panic("controller must be passed as a pointer")
	}
	if !composedOf(controller, Controller{}) {
		panic("controller must be composed of fuel.Controller")
	}

	// store it
	s.controllers = append(s.controllers, controller)
}

func (s *Server) loadEndpoints() {

	for _, controller := range s.controllers {

		// load endpoints::

		ctype := reflect.TypeOf(controller).Elem()
		cvalue := reflect.ValueOf(controller).Elem()

		// build controller fixture
		fixContr := func() Fixture {
			fldCont, _ := ctype.FieldByName("Controller")
			fixTag := newFixture(fldCont.Tag)
			fixCode := cvalue.FieldByName("Controller").FieldByName("Fixture").Interface().(Fixture)
			fixCode.Parent = &fixTag
			fixTag.Parent = &s.Fixture

			// if there is no root value set, then
			// use controller name (minus -controller) as Root
			if fixCode.getRoot() == "" {
				fixCode.Root = txt.CaseURL(ctype.Name())
				if strings.HasSuffix(fixCode.Root, "-controller") {
					fixCode.Root = fixCode.Root[0 : len(fixCode.Root)-len("-controller")]
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
					out.Route = txt.CaseSnake(fieldType.Name)
				}

				return out
			}()

			// build the endpoint, and store it in the server
			epoint := newEndpoint(fix, controller, fieldType, s)
			uniqURL := epoint.uniqueURL()
			if _, ok := s.endpoints[uniqURL]; ok {
				panic("cannot use same url again: " + uniqURL)
			}
			s.endpoints[uniqURL] = epoint
			fmt.Println(uniqURL)
		}
	}
}

func (s *Server) Run() {

	// load endpoints:
	// we do it at the end in the 'Run' step because
	// the user may add a cache later on (after calling AddController),
	// or the user may add middleware later on
	s.loadEndpoints()

	s.Addr = fmt.Sprintf(":%d", s.Port)
	s.Server.Handler = s.mux
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}