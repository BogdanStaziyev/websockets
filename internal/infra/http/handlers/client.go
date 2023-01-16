package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"myWebsockets/internal/domain"
	"net/http"
)

type Client struct {
	hub  *domain.Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn, hub *domain.Hub) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 1024),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (cli *Client) disconnected() {
	cli.hub.Unregister <- cli
	cli.conn.Close()
}

func (cli *Client) Socket(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	return nil
}
