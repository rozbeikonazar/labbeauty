package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cosmetcab.dp.ua/internal/assert"
	"github.com/gorilla/sessions"
)

// Mock handler for testing
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TestSecureHeaders tests the secureHeaders middleware
func TestSecureHeaders(t *testing.T) {
	app := &application{}
	handler := http.HandlerFunc(mockHandler)
	securedHandler := app.secureHeaders(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	securedHandler.ServeHTTP(rec, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":           "nosniff",
		"Strict-Transport-Security":        "max-age=31536000; includeSubDomains",
		"Access-Control-Allow-Origin":      "https://cosmetcab.dp.ua/",
		"Access-Control-Allow-Credentials": "true",
	}
	for key, value := range expectedHeaders {
		if got := rec.Header().Get(key); got != value {
			t.Errorf("Expected header %s, to be %s, but got %s", key, value, got)
		}
	}
}

// TestRecoverPanic tests the recoverPanic middleware
func TestRecoverPanic(t *testing.T) {
	app := &application{}
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Create a test handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	recoverHandler := app.recoverPanic(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	recoverHandler.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}

// TestRateLimit tests the rateLimit middleware
func TestRateLimit(t *testing.T) {
	app := &application{}
	app.config.limiter.burst = 1
	app.config.limiter.enabled = true
	app.config.limiter.rps = 1

	handler := http.HandlerFunc(mockHandler)
	rateLimitedHandler := app.rateLimit(handler)
	req := httptest.NewRequest("GET", "/", nil)
	rec1 := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()
	rec3 := httptest.NewRecorder()

	// First request, within rate limit
	rateLimitedHandler.ServeHTTP(rec1, req)
	assert.Equal(t, rec1.Code, http.StatusOK)

	// Second request, exceeds rate limit
	rateLimitedHandler.ServeHTTP(rec2, req)
	assert.Equal(t, rec2.Code, http.StatusTooManyRequests)

	// Third request from a different IP, within rate limit
	req.RemoteAddr = "192.168.0.1:5000"
	rateLimitedHandler.ServeHTTP(rec3, req)
	assert.Equal(t, rec3.Code, http.StatusOK)
}

// TestCheckAuth tests the checkAuth middleware
func TestCheckAuth(t *testing.T) {
	app := &application{
		sessionManager: sessions.NewCookieStore([]byte("test_token")),
	}

	handler := http.HandlerFunc(mockHandler)
	authHandler := app.checkAuth(handler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()
	rec3 := httptest.NewRecorder()
	session, _ := app.sessionManager.Get(req, "cookie-auth")

	// Unauthorized request
	session.Values["authenticated"] = false
	session.Save(req, rec)
	authHandler.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	// Authorized request
	session.Values["authenticated"] = true
	session.Save(req, rec2)
	authHandler.ServeHTTP(rec2, req)
	assert.Equal(t, rec2.Code, http.StatusOK)

	// Session not set
	authHandler.ServeHTTP(rec3, req)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)

}
