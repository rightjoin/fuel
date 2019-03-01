package fuel

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bugsnag/osext"
)

var dblSlash *regexp.Regexp

func cleanMultSlash(url string) string {
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

func serverInstance(run bool, conrollers ...serviceComposite) *Server {

	var server = NewServer()
	for _, c := range conrollers {
		server.AddService(c)
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

func readFile(path string) (string, error) {

	var absPath string
	var exists bool

	// if absolute path, try to find the file directly
	if strings.HasPrefix(path, "/") {
		absPath = path
		_, err := os.Stat(absPath)
		exists = err == nil
	}
	// try to locate it in working directory
	if !exists {
		wdir, ferr := os.Getwd()
		if ferr == nil {
			absPath = cleanMultSlash(wdir + "/" + path)
			_, err := os.Stat(absPath)
			exists = err == nil
		}
	}
	// try to locate it in executable directory
	if !exists {
		edir, ferr := osext.ExecutableFolder()
		if ferr == nil {
			absPath = cleanMultSlash(edir + "/" + path)
			_, err := os.Stat(absPath)
			exists = err == nil
		}
	}

	if !exists {
		return "", errors.New("File not found in working/exe path: " + path)
	}

	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
