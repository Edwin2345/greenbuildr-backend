package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := getEnv("PORT", "4000")
	clientSecret := getEnv("CLIENT_SECRET", "")
	apolloRouterURL := getEnv("APOLLO_ROUTER_URL", "http://localhost:4100")

	handler := NewHandler(clientSecret, apolloRouterURL)

	router := chi.NewRouter()
	router.Get("/health", healthHandler)
	router.Handle("/*", http.HandlerFunc(handler.Handle))

	log.Printf("apollo-auth proxy listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
