package session

import "net/http"

func PrepareRoutes(parent *http.ServeMux) {
	parent.Handle("POST /session/signup", http.HandlerFunc(handleSignUp))
	parent.HandleFunc("POST /session/login", handleLogIn)
	parent.HandleFunc("POST /session/logout", handleLogOut)
	parent.HandleFunc("POST /session/dropout", handleDropOut)
}
