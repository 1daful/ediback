// cors.go
package main

import (
	"github.com/rs/cors"
)

func SetupCORS(cfg Config) *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})
	return c
}
