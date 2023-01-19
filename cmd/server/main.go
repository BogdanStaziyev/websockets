package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"myWebsockets/config"
	"myWebsockets/config/constructor"
	"myWebsockets/internal/infra/http"
	"myWebsockets/internal/infra/http/routes"
)

func main() {
	var conf = config.GetConf()
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading ..env file")
	}

	server := http.NewServer()

	cont := constructor.New(*conf, *server)

	routes.Router(server, cont)

	err = server.Start()
	if err != nil {
		log.Fatal("Port already used")
	}
}
