package fuel

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/inconshreveable/log15"
)

func MiddlewareAccessLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info(r.RequestURI, "time", fmt.Sprintf("%.3fs", time.Now().Sub(start).Seconds()))
		})
	}
}

func MiddlewareAccessAndSlowLog(slowSeconds float64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			span := time.Now().Sub(start).Seconds()
			if span > slowSeconds {
				log.Warn(r.RequestURI, "time", fmt.Sprintf("%.3fs", span))
			} else {
				log.Info(r.RequestURI, "time", fmt.Sprintf("%.3fs", span))
			}
		})
	}
}
