package tests

import (
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type FileService struct {
	fuel.Service
	static fuel.GET `folder:"./sub/assets/"`
}

func TestStaticFileServer(t *testing.T) {
	server := fuel.NewServer()
	server.AddService(&FileService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

	web.Get("/file/static/abc.html").
		Expect(t).
		StatusOk().
		Done()

	web.Get("/file/static/css/def.css").
		Expect(t).
		StatusOk().
		Done()

	// TODO:
	// test that middleware chaining works on file server
}
