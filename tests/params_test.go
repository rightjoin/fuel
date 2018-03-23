package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type ParamController struct {
	fuel.Controller
	p1 fuel.GET `route:"p1/{string}/{int}/{uint}"`
}

func (s *ParamController) P1(str string, i int, u uint) string {
	return fmt.Sprintf("%s:%d:%d", str, i, u)
}

func TestParameters(t *testing.T) {
	server := fuel.NewServer()
	server.AddController(&ParamController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	web.Get("/param/p1/anything/1234/4321").
		Expect(t).
		Status(200).
		BodyEquals("anything:1234:4321").
		Done()

	web.Get("/param/p1/thing/-1234/4444").
		Expect(t).
		Status(200).
		BodyEquals("thing:-1234:4444").
		Done()
}
