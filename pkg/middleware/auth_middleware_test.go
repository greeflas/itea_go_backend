package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	header := http.Header{}
	header.Set("x-token", "test")

	req := &http.Request{
		Method: http.MethodPost,
		Header: header,
	}
	rw := httptest.NewRecorder()

	middleware := NewAuthMiddleware("test")
	wrappedHandler := middleware.Wrap(handler)

	wrappedHandler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusOK)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	header := http.Header{}
	header.Set("x-token", "invalid")

	req := &http.Request{
		Method: http.MethodPost,
		Header: header,
	}
	rw := httptest.NewRecorder()

	middleware := NewAuthMiddleware("test")
	wrappedHandler := middleware.Wrap(handler)

	wrappedHandler.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusUnauthorized)
	}
}
