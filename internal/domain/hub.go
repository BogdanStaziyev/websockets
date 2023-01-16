package domain

import (
	"github.com/go-redis/redis/v8"
	"myWebsockets/internal/infra/http/handlers"
)

type Hub struct {
	Rdb        *redis.Client
	Clients    map[*handlers.Client]bool
	Broadcast  chan []byte
	Register   chan *handlers.Client
	Unregister chan *handlers.Client
}
