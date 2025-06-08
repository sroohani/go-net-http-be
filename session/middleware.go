package session

import "net/http"

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || !isSessionValid(cookie.Value) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please login first"))
			return
		}
		next(w, r)
	}
}
