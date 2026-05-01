module github.com/greenbuildr/graphql-gateway

go 1.23

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/greenbuildr/shared v0.0.0
	go.uber.org/zap v1.27.0
)

replace github.com/greenbuildr/shared => ../shared
