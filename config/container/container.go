package container

import (
	"myWebsockets/config"
	"myWebsockets/internal/infra/http/handlers"
)

type Container struct {
	Services
	Handlers
	Middleware
}

type Services struct {
}

type Handlers struct {
	handlers.Client
}

type Middleware struct {
}

func New(conf config.Config) Container {

	clientHandler := handlers.NewClient()

	return Container{
		Services: Services{},
		Handlers: Handlers{
			clientHandler,
		},
		Middleware: Middleware{},
	}
}
