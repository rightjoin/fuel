package tests

import (
	"strconv"
	"testing"
	"time"

	"github.com/rightjoin/stag"
	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"

	"github.com/rightjoin/fuel"
)

type CacheController struct {
	fuel.Controller
	fewSeconds fuel.GET `cache:"store" ttl:"1s"`
}

func (s *CacheController) FewSeconds() string {
	time.Sleep(1 * time.Second)
	return "Slow"
}

func TestCaching(t *testing.T) {
	server := fuel.NewServer()
	server.DefineCache("store", stag.NewGoCache(5*time.Second))
	server.AddController(&CacheController{})
	port := runAsync(&server)

	var web = baloo.New("http://localhost:" + strconv.Itoa(port))

	var start time.Time
	var span time.Duration

	// first call should take > 1 sec
	start = time.Now()
	web.Get("/cache/few-seconds").
		Expect(t).
		Status(200).
		BodyEquals("Slow").
		Done()
	span = time.Now().Sub(start)
	assert.True(t, span >= (1*time.Second) && span < (2*time.Second))

	// next hit should be fast
	start = time.Now()
	web.Get("/cache/few-seconds").
		Expect(t).
		Status(200).
		BodyEquals("Slow").
		Done()
	span = time.Now().Sub(start)
	assert.True(t, span < (100*time.Millisecond))

	// expire the cache by waiting (1 sec)
	time.Sleep(1050 * time.Millisecond)

	// call should again take > 1 sec
	start = time.Now()
	web.Get("/cache/few-seconds").
		Expect(t).
		Status(200).
		BodyEquals("Slow").
		Done()
	span = time.Now().Sub(start)
	assert.True(t, span >= (1*time.Second) && span < (2*time.Second))
}
