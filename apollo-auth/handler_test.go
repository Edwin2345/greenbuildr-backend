package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func makeTestJWT(secret []byte, userID, email string) string {
	claims := &clientClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		panic(err)
	}
	return token
}

func TestHandle_StripsClientIdentityHeaders(t *testing.T) {
	var capturedHeaders http.Header
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
	}))
	defer backend.Close()

	handler := NewHandler("secret", backend.URL)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("X-User-ID", "spoofed-id")
	req.Header.Set("X-User-Email", "spoofed@evil.com")

	handler.Handle(httptest.NewRecorder(), req)

	if capturedHeaders.Get("X-User-ID") != "" {
		t.Error("X-User-ID should be stripped before forwarding")
	}
	if capturedHeaders.Get("X-User-Email") != "" {
		t.Error("X-User-Email should be stripped before forwarding")
	}
}

func TestHandle_InvalidJWT_Returns401(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not reach backend with invalid JWT")
	}))
	defer backend.Close()

	handler := NewHandler("secret", backend.URL)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}

	var body map[string]any
	json.NewDecoder(rr.Body).Decode(&body)
	errors, ok := body["errors"].([]any)
	if !ok || len(errors) == 0 {
		t.Fatal("expected errors array in response body")
	}
	firstError, ok := errors[0].(map[string]any)
	if !ok || firstError["message"] != "Unauthorized" {
		t.Errorf("expected message=Unauthorized, got %v", errors[0])
	}
}

func TestHandle_ValidJWT_InjectsUserHeaders(t *testing.T) {
	secret := []byte("test-secret")

	var capturedHeaders http.Header
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
	}))
	defer backend.Close()

	handler := NewHandler(string(secret), backend.URL)

	token := makeTestJWT(secret, "user-123", "user@example.com")
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	handler.Handle(httptest.NewRecorder(), req)

	if capturedHeaders.Get("X-User-ID") != "user-123" {
		t.Errorf("expected X-User-ID=user-123, got %q", capturedHeaders.Get("X-User-ID"))
	}
	if capturedHeaders.Get("X-User-Email") != "user@example.com" {
		t.Errorf("expected X-User-Email=user@example.com, got %q", capturedHeaders.Get("X-User-Email"))
	}
}

