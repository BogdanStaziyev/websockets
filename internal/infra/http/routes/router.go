package routes

import (
	"github.com/labstack/echo/v4/middleware"
	"myWebsockets/config/constructor"
	"myWebsockets/internal/infra/http"
)

func Router(s *http.Server, c constructor.Container) {
	e := s.Echo
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("ws", c.Handlers.Socket)
}
