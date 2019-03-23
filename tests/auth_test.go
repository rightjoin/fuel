package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type AuthService struct {
	fuel.Service

	locked fuel.GET `middleware:"basic"`
}

func (s *AuthService) Locked(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola")
}

func TestBasicAuth(t *testing.T) {
	server := fuel.NewServer()
	server.DefineMiddleware("basic", fuel.MidBasicAuth("abc", "def", "local"))
	server.AddService(&AuthService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

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
