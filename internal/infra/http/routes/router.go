package routes

import (
	"github.com/labstack/echo/v4/middleware"
	"myWebsockets/internal/infra/http"
	"myWebsockets/internal/infra/http/handlers"
)

func Router(s *http.Server) {
	webHandler := handlers.NewWebsocketConn(s)

	e := s.Echo
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("ws", webHandler.Socket)
}
