package main

import (
	"fmt"
	"os"
	"strconv"
)

type App struct {
	serverHost string
	serverPort int
	bcryptCost int
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

	bcryptCostStr := os.Getenv("BCRYPT_COST")
	if bcryptCostStr == "" {
		bcryptCostStr = "10"
	}
	bcryptCost, err := strconv.Atoi(bcryptCostStr)
	if err != nil {
		fmt.Println(err)
		bcryptCost = 10
	}
	a.bcryptCost = bcryptCost
}

func (a *App) ServerHost() string {
	return a.serverHost
}

func (a *App) ServerPort() int {
	return a.serverPort
}

func (a *App) BcryptCost() int {
	return a.bcryptCost
}
