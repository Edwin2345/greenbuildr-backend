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
	"github.com/jmoiron/sqlx"

	cryptoadapter "github.com/greenbuildr/auth-service/adapters/crypto"
	dbadapter "github.com/greenbuildr/auth-service/adapters/db"
	emailadapter "github.com/greenbuildr/auth-service/adapters/email"
	"github.com/greenbuildr/auth-service/domain/services"
	"github.com/greenbuildr/auth-service/graph"
)

func main() {
	//load env variables
	port := getEnv("PORT", "4001")
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root_password")
	dbName := getEnv("DB_NAME", "auth_db")
	resendKey := getEnv("RESEND_API_KEY", "")
	resendFrom := getEnv("RESEND_FROM", "onboarding@resend.dev")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	//set up db connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbName,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	//set up adapters
	repo := dbadapter.NewMySQLRepository(db)
	hasher := cryptoadapter.NewBcryptHasher()
	emailSender := emailadapter.NewResendSender(resendKey, resendFrom)

	//set up services
	registerSvc := services.NewRegisterService(repo, repo, hasher, emailSender, frontendURL)

	//Set Up GQL Server + Router
	resolver := &graph.Resolver{RegisterService: registerSvc}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	router := chi.NewRouter()
	router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("auth-service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
