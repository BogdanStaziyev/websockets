package http

import (
	"github.com/labstack/echo/v4/middleware"
	"myWebsockets/config/container"
)

func Router(s *Server, cont container.Container) {
	e := s.Echo
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("ws", cont.Client.Socket)
	e.Logger.Fatal(e.Start(":8080"))
}
