package app

import (
	"net/http"

	"github.com/rs/cors"
)

func handleCORS(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{`*`},
		AllowedMethods:   []string{`POST`, `GET`},
		AllowedHeaders:   []string{`*`},
		AllowCredentials: true,
		MaxAge:           600,
	}).Handler(h)
}
