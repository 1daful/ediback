// config.go
package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	AllowedOrigins []string
	AllowInsecure  bool
}

func LoadConfig() Config {
	port := getEnv("PORT", "2000")

	originsEnv := getEnv("ALLOWED_ORIGINS", "*")
	allowedOrigins := strings.Split(originsEnv, ",")

	allowInsecure, _ := strconv.ParseBool(getEnv("ALLOW_INSECURE", "false"))

	log.Printf("[config] Loaded | Port: %s | AllowedOrigins: %v | AllowInsecure: %v",
		port, allowedOrigins, allowInsecure)

	return Config{
		Port:           port,
		AllowedOrigins: allowedOrigins,
		AllowInsecure:  allowInsecure,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
