package tests

import (
	"strconv"
	"testing"

	baloo "gopkg.in/h2non/baloo.v3"

	"github.com/rightjoin/fuel"
)

type MockController struct {
	fuel.Controller
	yetToCode fuel.GET `stub:"sub/directory/stub_file.txt"`
}

func TestStubbing(t *testing.T) {
	server := fuel.NewServer()
	server.AddController(&MockController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	// first call should take > 1 sec
	web.Get("/mock/yet-to-code").
		Expect(t).
		Status(200).
		BodyEquals("file data").
		Done()
}
