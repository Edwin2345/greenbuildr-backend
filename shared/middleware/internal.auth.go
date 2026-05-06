package middleware

import (
	"context"
	"net/http"

	"github.com/greenbuildr/shared/contextkeys"
)

// InjectUserContext reads X-User-ID and X-User-Email headers injected by the
// apollo auth proxy and populates them into the request context. Returns 401 if either
// header is missing (unauthenticated request).
func InjectUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		email := r.Header.Get("X-User-Email")

		if userID == "" || email == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userID)
		ctx = context.WithValue(ctx, contextkeys.Email, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
