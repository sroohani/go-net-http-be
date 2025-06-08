package main

import (
	"os"
	"strconv"
)

type App struct {
	serverHost string
	serverPort int
}

func (a *App) Initialize() {
	a.serverHost = os.Getenv("SERVER_HOST")
	if a.serverHost == "" {
		a.serverHost = "localhost"
	}
	serverPortStr := os.Getenv("SERVER_HOST")
	if serverPortStr == "" {
		serverPortStr = "9876"
	}
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		serverPort = 9876
	}
	a.serverPort = serverPort
}

func (a *App) ServerHost() string {
	return a.serverHost
}

func (a *App) ServerPort() int {
	return a.serverPort
}
