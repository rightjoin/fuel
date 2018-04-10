package tests

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type AuthController struct {
	fuel.Controller

	locked fuel.GET `middleware:"basic"`
}

func (s *AuthController) Locked(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func TestBasicAuth(t *testing.T) {
	server := fuel.NewServer()
	server.DefineMiddleware("basic", fuel.MidBasicAuth("abc", "def", "local"))
	server.AddController(&AuthController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	web.Get("/auth/locked").
		Expect(t).
		Status(401).
		Done()

	r := web.Get("/auth/locked")
	r.Request.Context.Request.SetBasicAuth("abc", "def")
	r.Expect(t).
		Status(200).
		Done()

	r = web.Get("/auth/locked")
	r.Request.Context.Request.SetBasicAuth("abc", "wrong-password")
	r.Expect(t).
		Status(401).
		Done()

	r = web.Get("/auth/locked")
	r.Request.Context.Request.SetBasicAuth("wrong-username", "def")
	r.Expect(t).
		Status(401).
		Done()

}
