package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/greenbuildr/shared/contextkeys"
	"github.com/greenbuildr/shared/middleware"
)

const testSecret = "test-secret"

func makeToken(secret string, userID, email string, expiry time.Time) string {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   expiry.Unix(),
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	return token
}

// okHandler is a next handler that always returns 200, used in failure cases
// where we only care about the middleware rejecting the request.
func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestInternalAuth_MissingToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.DecodeInternalAuthToken(testSecret)(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInternalAuth_WrongSecret(t *testing.T) {
	token := makeToken("wrong-secret", "user-1", "user@example.com", time.Now().Add(time.Minute))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Token", token)
	rec := httptest.NewRecorder()

	middleware.DecodeInternalAuthToken(testSecret)(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInternalAuth_ExpiredToken(t *testing.T) {
	token := makeToken(testSecret, "user-1", "user@example.com", time.Now().Add(-time.Minute))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Token", token)
	rec := httptest.NewRecorder()

	middleware.DecodeInternalAuthToken(testSecret)(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInternalAuth_ValidToken(t *testing.T) {
	token := makeToken(testSecret, "user-1", "user@example.com", time.Now().Add(time.Minute))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Token", token)
	rec := httptest.NewRecorder()

	var capturedUserID, capturedEmail string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value(contextkeys.UserID).(string)
		capturedEmail = r.Context().Value(contextkeys.Email).(string)
		w.WriteHeader(http.StatusOK)
	})

	middleware.DecodeInternalAuthToken(testSecret)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if capturedUserID != "user-1" {
		t.Errorf("expected userID 'user-1', got %q", capturedUserID)
	}
	if capturedEmail != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got %q", capturedEmail)
	}
}
