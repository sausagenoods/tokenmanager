package tokenmanager

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

func getBearerToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// This middleware attempts to find the token in the store and
// saves it into request context.
func (t *TokenManager) TokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getBearerToken(r)

		data, err := t.Find(token)
		if err != nil {
			if errors.Is(err, ErrNotFound) || errors.Is(err, ErrExpired) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := t.SaveTokenToContext(r.Context(), data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (t *TokenManager) SaveTokenToContext(ctx context.Context, data *TokenInfo) context.Context {
	return context.WithValue(ctx, t.contextKey, data)
}

func (t *TokenManager) FromContext(ctx context.Context) *TokenInfo {
	c, ok := ctx.Value(t.contextKey).(*TokenInfo)
	if !ok {
		panic("tokenmanager: no token in context")
	}
	return c
}
