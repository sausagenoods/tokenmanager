package tokenmanager

import "net/http"

type tokenData {
	token string
	expiry time.Time
	data map[string]interface{}
}

func getBearerToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func (t *TokenManager) TokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getBearerToken(r)

		data, found, err := t.Find(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		//r.Context() = context.WithValue(ctx, , sd)
	}
}
