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
	fewSeconds fuel.GET `cache:"store" ttl:"250ms"`
}

func (s *CacheController) FewSeconds() string {
	time.Sleep(250 * time.Millisecond)
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

	// first call should take 250ms
	start = time.Now()
	web.Get("/cache/few-seconds").
		Expect(t).
		Status(200).
		BodyEquals("Slow").
		Done()
	span = time.Now().Sub(start)
	assert.True(t, span >= (250*time.Millisecond) && span < (300*time.Millisecond))

	// next hit should be fast
	start = time.Now()
	web.Get("/cache/few-seconds").
		Expect(t).
		Status(200).
		BodyEquals("Slow").
		Done()
	span = time.Now().Sub(start)
	assert.True(t, span < (100*time.Millisecond))

	// expire the cache by waiting for ~250ms
	time.Sleep(300 * time.Millisecond)

	// call should again take 250ms
	start = time.Now()
	web.Get("/cache/few-seconds").
		Expect(t).
		Status(200).
		BodyEquals("Slow").
		Done()
	span = time.Now().Sub(start)
	assert.True(t, span >= (250*time.Millisecond) && span < (300*time.Millisecond))
}
