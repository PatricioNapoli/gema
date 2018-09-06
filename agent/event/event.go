package event

import (
	"github.com/go-redis/redis"
	"gema/agent/database"
	"gema/agent/utils"
	"fmt"
)

type Attributes struct {
	ServiceName string `json:"com.docker.swarm.service.name"`
	Expose string `json:"gema.expose"`
	Port string `json:"gema.port"`
}

type Actor struct {
	Attributes Attributes
}

type Event struct {
	Status string `json:"status"`
	Actor Actor
}

type Handler struct {
	Database *redis.Client
}

func New() *Handler{
	return &Handler{
		Database: database.New(),
	}
}

func (h *Handler) HandleEvent(ev *Event) {
	if ev.Status != "create" {
		return
	}

	buf, count := utils.ReadFile(fmt.Sprintf("/config/%s", ev.Actor.Attributes.ConfigFile))

	err := h.Database.Set(ev.Actor.Attributes.ServiceName, string(buf[:count]), 0).Err()
	utils.Handle(err)

	println("Updated service routing for service: " + ev.Actor.Attributes.ServiceName)
}