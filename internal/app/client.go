package app

import (
	"encoding/json"
	"errors"
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

type ClientService interface {
	ProcessEvents(rawMessage []byte, client *domain.Client) error
	ReadPump(client *domain.Client)
	WritePump(client *domain.Client)
}

type clientService struct {
	e EventService
}

func NewClientService(event EventService) ClientService {
	return clientService{e: event}
}

func (c clientService) ProcessEvents(rawMessage []byte, client *domain.Client) error {
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
		if err = c.e.SendMessageToAll(client, message); err != nil {
			return err
		}
	case domain.ActionSendPrivate:
		var message domain.SendMessageToOne
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println(err)
			return err
		}
		if err = c.e.SendMessageToOne(client, message); err != nil {
			return err
		}
	}

	return err
}

// ReadPump pumps events from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c clientService) ReadPump(client *domain.Client) {

	defer func() {
		client.Disconnect()
	}()
	client.Conn.SetReadLimit(maxMessageSize)
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	client.Conn.SetPongHandler(func(string) error {
		//log.Println("pong")
		return client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("error: %v", err)
			break
		}
		err = c.ProcessEvents(message, client)
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
func (c clientService) WritePump(client *domain.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			//todo ???
			err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				if err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
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
			err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err = client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("write msg: ", err)
				return
			}
		}
	}
}
