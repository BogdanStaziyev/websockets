package http

import (
	"github.com/labstack/echo/v4"
	"myWebsockets/internal/app"
)

type Server struct {
	Echo *echo.Echo
	Hub  *app.Hub
}

func NewServer() *Server {
	s := &Server{
		Echo: echo.New(),
		Hub:  app.NewHub(),
	}
	go s.Hub.Run()
	return s
}

func (s *Server) Start() error {
	return s.Echo.Start(":8080")
}
