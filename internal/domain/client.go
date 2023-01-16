package domain

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
	hub  *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		hub:  hub,
	}
}
