package main

import (
	"log"
	"myWebsockets/internal/infra/http"
	"myWebsockets/internal/infra/http/routes"
)

func main() {
	//var conf = config.GetConf()
	//err := godotenv.Load()
	//if err != nil {
	//	fmt.Println("Error loading ..env file")
	//}
	//cont := container.New(*conf)

	server := http.NewServer()

	routes.Router(server)

	err := server.Start()
	if err != nil {
		log.Fatal("Port already used")
	}
}
