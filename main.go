// main.go
package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
)

func main() {
	_ = godotenv.Load()

	cfg := LoadConfig()

	r := mux.NewRouter()

	// Proxy route
	r.Handle("/api/v1/proxy", ProxyHandler(cfg)).Methods(http.MethodPost)

	// Optional admin routes
	r.HandleFunc("/admin", AdminDashboard).Methods(http.MethodGet)
	r.HandleFunc("/admin/run-task", AdminRunTask).Methods(http.MethodPost)

	handler := SetupCORS(cfg).Handler(r)

	log.Println("[server] Listening on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
