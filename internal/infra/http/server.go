package http

import (
	"github.com/labstack/echo/v4"
	"myWebsockets/internal/domain"
)

type Server struct {
	Echo *echo.Echo
	Hub  *domain.Hub
}

func NewServer() *Server {
	s := &Server{
		Echo: echo.New(),
		Hub:  domain.NewHub(),
	}
	go s.Hub.Run()
	return s
}

func (s *Server) Start() error {
	return s.Echo.Start(":8080")
}
