package fuel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanURL(t *testing.T) {
	var badURL = "//abc//def//ghi//"
	assert.Equal(t, "/abc/def/ghi/", cleanUrl(badURL))
}
