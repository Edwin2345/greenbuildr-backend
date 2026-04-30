module github.com/greenbuildr/auth-service

go 1.25.0

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	golang.org/x/crypto v0.50.0
)

require github.com/resend/resend-go/v2 v2.28.0 // indirect

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/jmoiron/sqlx v1.4.0
)

replace github.com/greenbuildr/shared => ../shared
