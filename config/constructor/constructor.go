package constructor

import (
	"myWebsockets/config"
	"myWebsockets/internal/app"
	"myWebsockets/internal/infra/http"
	"myWebsockets/internal/infra/http/handlers"
)

type Container struct {
	Services
	Handlers
	Middleware
}

type Services struct {
	app.EventService
	app.ClientService
}

type Handlers struct {
	handlers.WebsocketConn
}

type Middleware struct {
}

func New(conf config.Config, s http.Server) Container {

	eventService := app.NewEventService()

	clientService := app.NewClientService(eventService)
	clientHandler := handlers.NewWebsocketConn(&s, clientService)

	return Container{
		Services: Services{
			eventService,
			clientService,
		},
		Handlers: Handlers{
			clientHandler,
		},
		Middleware: Middleware{},
	}
}
