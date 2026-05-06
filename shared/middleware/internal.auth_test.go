package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greenbuildr/shared/contextkeys"
	"github.com/greenbuildr/shared/middleware"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestInjectUserContext_MissingHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.InjectUserContext(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInjectUserContext_MissingEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-1")
	rec := httptest.NewRecorder()

	middleware.InjectUserContext(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInjectUserContext_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-1")
	req.Header.Set("X-User-Email", "user@example.com")
	rec := httptest.NewRecorder()

	var capturedUserID, capturedEmail string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value(contextkeys.UserID).(string)
		capturedEmail = r.Context().Value(contextkeys.Email).(string)
		w.WriteHeader(http.StatusOK)
	})

	middleware.InjectUserContext(next).ServeHTTP(rec, req)

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
