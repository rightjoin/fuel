package tests

import (
	"strconv"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type MvcController struct {
	fuel.Controller
	hello     fuel.GET
	override  fuel.GET
	override2 fuel.GET
}

func (s *MvcController) Hello() fuel.View {
	return fuel.View{}
}

func (s *MvcController) Override() fuel.View {
	return fuel.View{
		Layout: "layout1",
	}
}

func (s *MvcController) Override2() fuel.View {
	return fuel.View{
		View: "hola",
	}
}

// func (s *MvcController) Partials() fuel.View {
// 	return fuel.View{}
// }

func TestViewDefaultOption(t *testing.T) {
	server := fuel.NewServer()
	server.AddController(&MvcController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	web.Get("/mvc/hello").
		Expect(t).
		Status(200).
		BodyEquals("View").
		Done()
}

func TestCustomViewFolder(t *testing.T) {
	server := fuel.NewServer()
	server.MvcOptions.Views = "templates/custom_folder"
	server.MvcOptions.Layout = "layout"
	server.AddController(&MvcController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	web.Get("/mvc/hello").
		Expect(t).
		Status(200).
		BodyEquals("Start View Close").
		Done()
}

func TestLayoutOverrideFolder(t *testing.T) {
	server := fuel.NewServer()
	server.MvcOptions.Views = "templates/override_layout"
	//server.MvcOptions.Layout = "layout"
	server.AddController(&MvcController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	web.Get("/mvc/override").
		Expect(t).
		Status(200).
		BodyEquals("Start1 View1 Close1").
		Done()
}

func TestViewOverride(t *testing.T) {
	server := fuel.NewServer()
	server.MvcOptions.Views = "templates/override_view"
	server.MvcOptions.Layout = "layout"
	server.AddController(&MvcController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	web.Get("/mvc/override2").
		Expect(t).
		Status(200).
		BodyEquals("Start View2 Close").
		Done()
}

// TODO: partials
// func TestPartials(t *testing.T) {
// 	server := fuel.NewServer()
// 	server.MvcOptions.Views = "templates/partials"
// 	server.MvcOptions.Layout = "layout"
// 	server.AddController(&MvcController{})
// 	port := runAsync(&server)

// 	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

// 	web.Get("/mvc/partials").
// 		Expect(t).
// 		Status(200).
// 		BodyEquals("Start View2 Close").
// 		Done()
// }
