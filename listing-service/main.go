package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/greenbuildr/listing-service/graph"
)

const defaultPort = "4002"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// Initialize resolver with database connection
	// TODO: Add database connection
	resolver := &graph.Resolver{}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// GraphQL endpoint
	router.Handle("/graphql", srv)

	// GraphQL Playground for development
	router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))

	log.Printf("Listing Service listening on http://localhost:%s/", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
