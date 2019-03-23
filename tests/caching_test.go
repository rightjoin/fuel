package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"

	"github.com/rightjoin/fuel"
	"github.com/rightjoin/stak"
)

type CacheService struct {
	fuel.Service
	fewSeconds fuel.GET `cache:"store" ttl:"250ms"`
}

func (s *CacheService) FewSeconds() string {
	time.Sleep(250 * time.Millisecond)
	return "Slow"
}

func TestCaching(t *testing.T) {
	server := fuel.NewServer()
	server.DefineCache("store", stak.NewGoCache(5*time.Second))
	server.AddService(&CacheService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

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
