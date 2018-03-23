package tests

import (
	"time"

	"github.com/rightjoin/fuel"
)

var port = 8080

func runAsync(s *fuel.Server) int {
	port++
	s.Port = port
	go s.Run()

	// wait some time for the server to fire up
	time.Sleep(50 * time.Millisecond)

	return port
}
