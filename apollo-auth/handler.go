package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	clientSecret []byte
	proxy        *httputil.ReverseProxy
}

func NewHandler(clientSecret, apolloRouterURL string) *Handler {
	target, _ := url.Parse(apolloRouterURL)
	return &Handler{
		clientSecret: []byte(clientSecret),
		proxy:        httputil.NewSingleHostReverseProxy(target),
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Strip any client-sent identity headers before we inject our own
	r.Header.Del("X-User-ID")
	r.Header.Del("X-User-Email")

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		//parse and validate client jwt
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, email, err := h.validateClientJWT(tokenStr)
		if err != nil {
			log.Printf("apollo-auth-proxy: invalid client JWT: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{
				"errors": []map[string]string{{"message": "Unauthorized"}},
			})
			return
		}

		//add trusted headers
		r.Header.Set("X-User-ID", userID)
		r.Header.Set("X-User-Email", email)
		log.Printf("apollo-auth-proxy: forwarding request for userID=%s", userID)
	}

	//proxy request to appollo router
	h.proxy.ServeHTTP(w, r)
}

type clientClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (h *Handler) validateClientJWT(tokenStr string) (string, string, error) {
	claims := &clientClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return h.clientSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}
	return claims.Subject, claims.Email, nil
}
