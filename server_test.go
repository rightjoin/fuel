package fuel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testControllerA struct {
}

func (t testControllerA) BeginRequest() {

}

func (t testControllerA) EndRequest() {

}

func TestAddControllerNotComposedOfController(t *testing.T) {

	s := NewServer()
	assert.Panics(t, func() {
		s.AddController(&testControllerA{})
	})
}

func TestAddControllerAddress(t *testing.T) {

	type testControllerB struct {
		Controller
	}

	s := NewServer()
	assert.Panics(t, func() {
		s.AddController(testControllerB{})
	})
	assert.NotPanics(t, func() {
		s.AddController(&testControllerB{})
	})
}

type DittooController struct {
	Controller
	aGet GET `route:"a-url"`
	bGet GET `route:"a-url"`
}

func (me DittooController) AGet() string { return "" }
func (me DittooController) BGet() string { return "" }

func TestDittooURLs(t *testing.T) {
	s := NewServer()
	s.AddController(&DittooController{})
	assert.Panics(t, func() {
		s.loadEndpoints()
	})
}
