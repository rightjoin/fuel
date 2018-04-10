package fuel

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	log "github.com/inconshreveable/log15"
)

func MidAccessLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info(r.RequestURI, "time", fmt.Sprintf("%.3fs", time.Now().Sub(start).Seconds()))
		})
	}
}

func MidAccessAndSlowLog(slowSeconds float64) func(http.Handler) http.Handler {
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

func MidBasicAuth(username, password, realm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// https: //stackoverflow.com/questions/21936332/idiomatic-way-of-requiring-http-basic-auth-in-go
			user, pass, ok := r.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
