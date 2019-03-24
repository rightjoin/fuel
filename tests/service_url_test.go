package tests

import (
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type UrlService struct {
	fuel.Service
	sayHello fuel.GET
	sayHola  fuel.GET
	postOne  fuel.POST `route:"one"`
	postTwo  fuel.POST `route:"-"`
}

func (s *UrlService) SayHello() string {
	return "Hello World"
}

func (s *UrlService) SayHola() string {
	return "Hello World"
}

func (s *UrlService) PostOne() string {
	return "One"
}

func (s *UrlService) PostTwo() string {
	return "Two"
}

func TestUrlStructure(t *testing.T) {
	server := fuel.NewServer()
	server.AddService(&UrlService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

	web.Get("/url/say-hello").
		Expect(t).
		Status(200).
		Done()

	web.Get("/url/say-hola").
		Expect(t).
		Status(200).
		Done()

	web.Post("/url/one").
		Expect(t).
		StatusOk().
		Done()

	web.Post("/url").
		Expect(t).
		StatusOk().
		Done()

}
