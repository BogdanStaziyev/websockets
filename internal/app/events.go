package app

import (
	"encoding/json"
	"log"
	"myWebsockets/internal/domain"
)

type EventService interface {
	SendMessageToAll(c *domain.Client, message domain.SendMessageToAll) error
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
