module github.com/greenbuildr/auth-service

go 1.23

require (
	github.com/99designs/gqlgen v0.17.49
	github.com/go-chi/chi/v5 v5.0.12
	github.com/go-sql-driver/mysql v1.7.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/greenbuildr/shared v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.16
	golang.org/x/crypto v0.24.0
)

replace github.com/greenbuildr/shared => ../shared
