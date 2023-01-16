package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"myWebsockets/internal/infra/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	http.NewServer()
}
