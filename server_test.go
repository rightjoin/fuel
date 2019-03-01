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
