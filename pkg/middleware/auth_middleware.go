package middleware

import "net/http"

type AuthMiddleware struct {
	token string
}

func NewAuthMiddleware(token string) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPatch, http.MethodDelete:
			if r.Header.Get("x-token") != m.token {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
