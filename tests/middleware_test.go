package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type MiddlewareService struct {
	fuel.Service
	m1 fuel.GET `route:"m1/{input}"`
	m2 fuel.GET `route:"m2/{input}" middle:"b"`
	m3 fuel.GET `route:"m3/{input}" middleware:"a,b"`
}

func (s *MiddlewareService) M1(inp string) string {
	return "M1"
}

func (s *MiddlewareService) M2(inp string) string {
	return "M2"
}

func (s *MiddlewareService) M3(inp string) string {
	return "M3"
}

// return 'a' if request URI ends with a
func midA() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.RequestURI, "a") {
				fmt.Fprintf(w, "a")
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// return 'b' if request URI ends with b
func midB() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.RequestURI, "b") {
				fmt.Fprintf(w, "b")
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	server := fuel.NewServer()
	server.AddService(&MiddlewareService{})
	server.DefineMiddleware("a", midA())
	server.DefineMiddleware("b", midB())
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

	// no middleware
	web.Get("/middleware/m1/anything").
		Expect(t).
		Status(200).
		BodyEquals("M1").
		Done()

	// middleware b on m2
	web.Get("/middleware/m2/anything").
		Expect(t).
		Status(200).
		BodyEquals("M2").
		Done()
	web.Get("/middleware/m2/bbb").
		Expect(t).
		Status(200).
		BodyEquals("b").
		Done()

	// middleware a, b on m3
	web.Get("/middleware/m3/anything").
		Expect(t).
		Status(200).
		BodyEquals("M3").
		Done()
	web.Get("/middleware/m3/bbb").
		Expect(t).
		Status(200).
		BodyEquals("b").
		Done()
	web.Get("/middleware/m3/aaa").
		Expect(t).
		Status(200).
		BodyEquals("a").
		Done()
}
