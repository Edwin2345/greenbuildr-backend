package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/greenbuildr/shared/contextkeys"
)

type internalClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// validates the internal JWT injected by the gateway and populates userID and email into the request context.
func DecodeInternalAuthToken(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := strings.TrimPrefix(r.Header.Get("X-Internal-Token"), "Bearer ")
			if tokenStr == "" {
				http.Error(w, "missing internal token", http.StatusUnauthorized)
				return
			}

			claims := &internalClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid internal token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserID, claims.Subject)
			ctx = context.WithValue(ctx, contextkeys.Email, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
