// cors.go
package main

import (
	"net/http"

	"github.com/rs/cors"
)

func SetupCORS(cfg Config) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}
