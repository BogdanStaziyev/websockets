package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	for {
		// Read
		_, msg, err := conn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
