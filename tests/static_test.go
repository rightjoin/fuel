package tests

import (
	"strconv"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type FileController struct {
	fuel.Controller
	static fuel.GET `folder:"./sub/assets/"`
}

func TestStaticFileServer(t *testing.T) {
	server := fuel.NewServer()
	server.AddController(&FileController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

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
