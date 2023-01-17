package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"myWebsockets/internal/app"
	s "myWebsockets/internal/infra/http"
	"net/http"
)

type WebsocketConn struct {
	server *s.Server
}

func NewWebsocketConn(s *s.Server) *WebsocketConn {
	return &WebsocketConn{
		server: s,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (cli *WebsocketConn) Socket(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	client := app.NewClient(conn, cli.server.Hub)

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
	return err
}
