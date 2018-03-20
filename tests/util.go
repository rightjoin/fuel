package tests

import (
	"time"

	"github.com/rightjoin/fuel"
)

var port = 8421

func asyncRun(s *fuel.Server) int {
	port++
	s.Port = port
	go s.Run()

	// wait some time for the server to fire up
	time.Sleep(50 * time.Millisecond)

	// TODO:
	// check after 1 second. if the server is not up, then panic

	return port
}
