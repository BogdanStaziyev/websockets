package http

import (
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
}

func NewServer() *Server {
	return &Server{Echo: echo.New()}
}

func (s Server) Start() error {
	return s.Echo.Start(":8080")
}
