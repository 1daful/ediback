// config.go
package main

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Port           string
	AllowedHosts   []string
	AllowedOrigins []string
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2000"
	}

	allowedHosts := strings.Split(os.Getenv("ALLOWED_HOSTS"), ",")
	if len(allowedHosts) == 0 || allowedHosts[0] == "" {
		allowedHosts = []string{"localhost", "127.0.0.1"}
	}

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || allowedOrigins[0] == "" {
		allowedOrigins = []string{"*"}
	}

	log.Println("[proxy-server] Config Loaded")
	log.Println("Port:", port)
	log.Println("Allowed Hosts:", allowedHosts)
	log.Println("Allowed Origins:", allowedOrigins)

	return Config{
		Port:           port,
		AllowedHosts:   allowedHosts,
		AllowedOrigins: allowedOrigins,
	}
}
