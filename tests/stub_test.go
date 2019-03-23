package tests

import (
	"testing"

	baloo "gopkg.in/h2non/baloo.v3"

	"github.com/rightjoin/fuel"
)

type MockService struct {
	fuel.Service
	yetToCode fuel.GET `stub:"sub/directory/stub_file.txt"`
}

func TestStubbing(t *testing.T) {
	server := fuel.NewServer()
	server.AddService(&MockService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

	// first call should take > 1 sec
	web.Get("/mock/yet-to-code").
		Expect(t).
		Status(200).
		BodyEquals("file data").
		Done()
}
