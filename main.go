package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sroohani/go-net-http-be/jwt"
	"github.com/sroohani/go-net-http-be/session"
)

type DefaultResponse struct {
	Message string   `json:"message"`
	Routes  []string `json:"routes"`
}

type HttpHandler struct {
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(DefaultResponse{
		Message: "Please send POST requests to one of the available routes.",
		Routes: []string{
			"/session/signup", "/session/login", "/session/logout",
			"/jwt/signup", "/jwt/login", "/jwt/logout",
		}})
	if err != nil {
		return
	}
	w.Write(res)
}

func main() {
	_ = godotenv.Load()

	a := &App{}
	a.Initialize()
	session.GenerateSecret()
	session.SetBcryptCost(a.BcryptCost())

	router := http.NewServeMux()
	router.Handle("GET /{$}", HttpHandler{})
	session.PrepareRoutes(router)
	jwt.PrepareRoutes(router)
	router.HandleFunc("/{p...}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		res, err := json.Marshal(DefaultResponse{Message: "Bad request!"})
		if err != nil {
			return
		}
		w.Write(res)
	})
	server := http.Server{
		Addr:    fmt.Sprintf("%v:%v", a.ServerHost(), a.ServerPort()),
		Handler: router,
	}

	server.ListenAndServe()

}
