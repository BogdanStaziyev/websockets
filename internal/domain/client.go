package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Conn: conn,
		Hub:  hub,
		Send: make(chan []byte, 1024),
	}
}

func (c *Client) disconnect() {
	c.Hub.Unregister <- c
	c.Conn.Close()
}

type base struct {
	Name string `json:"name"`
}

func (c *Client) processEvents(rawMessage []byte) error {
	var baseMessage base
	err := json.Unmarshal(rawMessage, &baseMessage)
	if err != nil {
		return err
	}

	if baseMessage.Name == "" {
		return errors.New("error deserializing message")
	}
	for clients := range c.Hub.Clients {
		if clients != nil {
			err = clients.Conn.WriteJSON(baseMessage)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				err = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error: %v", err)
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.disconnect()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	c.Conn.SetPongHandler(func(string) error { _ = c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("error: %v", err)
			break
		}
		err = c.processEvents(message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
