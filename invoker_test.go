package fuel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
}

func (s testStruct) GetOne() string {
	return "one"
}

func (s *testStruct) GetTwo() string {
	return "two"
}

func (s testStruct) GetThree(i int, str string) (float64, string) {
	return 0.5 * float64(i), str + str
}

func TestInvocationOnStruct(t *testing.T) {

	var s testStruct

	i1 := newInvoker(&s, "GetOne")
	assert.Equal(t, "one", i1.invoke()[0].String())

	i2 := newInvoker(s, "GetOne")
	assert.Equal(t, "one", i2.invoke()[0].String())

	i3 := newInvoker(&s, "GetTwo")
	assert.Equal(t, "two", i3.invoke()[0].String())

	assert.Panics(t, func() {
		i4 := newInvoker(s, "GetTwo")
		i4.invoke()
	})
}

func TestInputsOutputs(t *testing.T) {

	var s testStruct

	i1 := newInvoker(&s, "GetThree")
	assert.Equal(t, 5.0, i1.invoke(10, "abc")[0].Float())
	assert.Equal(t, "abab", i1.invoke(10, "ab")[1].String())
}
