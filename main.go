// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := LoadConfig()
	handler := SetupCORS(cfg).Handler(http.HandlerFunc(ProxyHandler(cfg)))

	log.Println("[proxy-server] Running on port:", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
