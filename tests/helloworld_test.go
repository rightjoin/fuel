package tests

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type HelloWorldController struct {
	fuel.Controller
	sayHello fuel.GET
	sayHola  fuel.GET
}

func (s *HelloWorldController) SayHello() string {
	return "Hello World"
}

func (s *HelloWorldController) SayHola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func TestHelloWorld(t *testing.T) {
	server := fuel.NewServer()
	server.AddController(&HelloWorldController{})
	port := asyncRun(&server)
	defer server.Close()

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

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
