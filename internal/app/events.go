package app

import (
	"encoding/json"
	"log"
	"myWebsockets/internal/domain"
)

type EventService interface {
	SendMessageToAll(c *domain.Client, message domain.SendMessageToAll) error
	SendMessageToOne(c *domain.Client, message domain.SendMessageToOne) error
}

type eventService struct {
}

func NewEventService() EventService {
	return eventService{}
}

func (e eventService) SendMessageToAll(c *domain.Client, message domain.SendMessageToAll) error {
	byt, err := json.Marshal(message.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	c.Hub.Broadcast <- byt
	return nil
}

func (e eventService) SendMessageToOne(c *domain.Client, message domain.SendMessageToOne) error {
	byt, err := json.Marshal(message.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	//todo find in db after register
	for ind := range c.Hub.Clients {
		if ind.ID == message.UserID {
			ind.Send <- byt
			return nil
		}
	}
	c.Send <- []byte("client does not exist")
	return nil
}
