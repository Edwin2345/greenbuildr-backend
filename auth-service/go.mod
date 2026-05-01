module github.com/greenbuildr/auth-service

go 1.25.0

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.50.0
)

require (
	github.com/99designs/gqlgen v0.17.90
	github.com/go-chi/chi/v5 v5.2.5
	github.com/greenbuildr/shared v0.0.0-00010101000000-000000000000
	github.com/resend/resend-go/v2 v2.28.0
	github.com/vektah/gqlparser/v2 v2.5.33
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/jmoiron/sqlx v1.4.0
)

replace github.com/greenbuildr/shared => ../shared
