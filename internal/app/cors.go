package app

import (
	"net/http"

	"github.com/rs/cors"
)

func handleCORS(rw http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodOptions {
		return false
	}
	if r.Header.Get("Access-Control-Request-Method") == "" {
		return false
	}

	stubH := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "", http.StatusServiceUnavailable)
	})

	(cors.New(cors.Options{
		AllowedOrigins:   []string{`*`},
		AllowedMethods:   []string{`POST`, `GET`},
		AllowedHeaders:   []string{`*`},
		AllowCredentials: true,
		MaxAge:           600,
	})).Handler(stubH).ServeHTTP(rw, r)
	return true
}
