package fuel

import (
	"regexp"
	"strings"
	"time"
)

var dblSlash *regexp.Regexp

func cleanUrl(url string) string {
	if dblSlash == nil {
		dblSlash, _ = regexp.Compile("[\\/]+")
	}

	return dblSlash.ReplaceAllString(url, "/")
}

var muxParse *regexp.Regexp

func extractMuxVars(url string) []string {

	if muxParse == nil {
		muxParse, _ = regexp.Compile(`{[^/]+}`)
	}

	matches := muxParse.FindAllString(url, -1)
	var colonPos int
	for i, m := range matches {
		m = m[1 : len(m)-1] // drop { and }
		colonPos = strings.Index(m, ":")
		if colonPos > 0 {
			m = m[0:colonPos]
		}
		matches[i] = m
	}

	return matches
}

func acceptableOutput(typeSym string) bool {

	// adjust for pointer
	if typeSym[0:1] == "*" {
		typeSym = typeSym[1:]
	}

	switch {
	case typeSym == "string":
	case typeSym == "map":
	case strings.HasPrefix(typeSym, "st:"):
	case strings.HasPrefix(typeSym, "sl:"):
	case strings.HasPrefix(typeSym, "i:"):
	default:
		return false
	}
	return true
}

var portTesting = 9900

func serverInstance(run bool, conrollers ...service) *Server {

	var server = NewServer()
	for _, c := range conrollers {
		server.AddController(c)
	}

	// TODO: lock
	portTesting++
	server.Port = portTesting

	if run {
		go server.Run()
		time.Sleep(time.Millisecond * 50)
	}

	return &server
}
