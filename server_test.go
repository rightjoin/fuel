package fuel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type serviceA struct {
}

func (t serviceA) BeginRequest() {

}

func (t serviceA) EndRequest() {

}

func TestAddServiceNotComposedOfFuelService(t *testing.T) {

	s := NewServer()
	assert.Panics(t, func() {
		s.AddService(&serviceA{})
	})
}

func TestAddServiceAddress(t *testing.T) {

	type serviceB struct {
		Service
	}

	s := NewServer()
	assert.Panics(t, func() {
		s.AddService(serviceB{})
	})
	assert.NotPanics(t, func() {
		s.AddService(&serviceB{})
	})
}

type serviceC struct {
	Service
	aGet GET `route:"a-url"`
	bGet GET `route:"a-url"`
}

func (me serviceC) AGet() string { return "" }
func (me serviceC) BGet() string { return "" }

func TestDittooURLs(t *testing.T) {
	s := NewServer()
	s.AddService(&serviceC{})
	assert.Panics(t, func() {
		s.loadEndpoints()
	})
}

type serviceD struct {
	Service
	one   POST `route:"one"`
	two   POST `route:"two"`
	three POST `route:"three/{a}"`
	four  POST `route:"four/{a}"`
}

type LoadMe struct {
	Field string `json:"field"`
}

func (me serviceD) One(s LoadMe) string                  { return "im-one" }
func (me serviceD) Two(s LoadMe, ad Aide) string         { return "im-two" }
func (me serviceD) Three(a int, s LoadMe) string         { return "im-three" }
func (me serviceD) Four(a int, s LoadMe, ad Aide) string { return "im-four" }

func TestStructLoading(t *testing.T) {
	s := NewServer()
	s.AddService(&serviceD{})
	s.RunTestInstance()
}
