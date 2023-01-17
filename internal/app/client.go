package app

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"myWebsockets/internal/domain"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
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

func (c *Client) disconnect() {
	c.Hub.Unregister <- c
	c.Conn.Close()
}

func (c *Client) processEvents(rawMessage []byte) error {
	var baseMessage domain.Base
	err := json.Unmarshal(rawMessage, &baseMessage)
	if err != nil {
		return err
	}

	if baseMessage.Action == "" {
		return errors.New("error deserializing message")
	}

	switch baseMessage.Action {
	case domain.ActionSandMessage:
		var message domain.SendMessageToAll
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println(err)
			return err
		}
		byt, err := json.Marshal(message.Message)
		if err != nil {
			log.Println(err)
			return err
		}
		c.Hub.Broadcast <- byt
	}

	return err
}

// ReadPump pumps events from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.disconnect()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		//log.Println("pong")
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})
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
			log.Printf("error: %v", err)
		}
	}

}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			//todo ???
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Println("connection closed error: ", err)
				return
			}
			if err := w.Close(); err != nil {
				log.Println("connection closed error: ", err)
				return
			}
		case <-ticker.C:
			//log.Printf("ping")
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("write msg: ", err)
				return
			}
		}
	}
}
