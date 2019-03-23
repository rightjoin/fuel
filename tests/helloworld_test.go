package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type HelloWorldService struct {
	fuel.Service
	sayHello fuel.GET
	sayHola  fuel.GET
}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

func (s *HelloWorldService) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func TestHelloWorld(t *testing.T) {
	server := fuel.NewServer()
	server.AddService(&HelloWorldService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

	web.Get("/hello-world/say-hello").
		Expect(t).
		Status(200).
		Header("Content-Type", "text/plain").
		BodyEquals("Hello World").
		Done()

	web.Get("/hello-world/say-hola").
		Expect(t).
		Status(200).
		BodyEquals("Hola").
		Done()
}
