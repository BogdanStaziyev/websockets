package domain

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Name string
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Conn: conn,
		Hub:  hub,
		Send: make(chan []byte, 1024),
		ID:   uuid.New().String(),
	}
}

func (c *Client) Disconnect() {
	c.Hub.Unregister <- c
	c.Conn.Close()
}
