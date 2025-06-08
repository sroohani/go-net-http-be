package jwt

import "net/http"

func PrepareRoutes(parent *http.ServeMux) {
	parent.Handle("POST /jwt/signup", http.HandlerFunc(handleSignUp))
	parent.HandleFunc("POST /jwt/login", handleLogIn)
	parent.HandleFunc("POST /jwt/logout", handleLogOut)
}
