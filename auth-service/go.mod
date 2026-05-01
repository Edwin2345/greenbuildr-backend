module github.com/greenbuildr/auth-service

go 1.25.0

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.50.0
)

require github.com/resend/resend-go/v2 v2.28.0

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
