package event

import (
	"strings"
	"fmt"
	"gema/agent/database"
	"gema/agent/utils"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

type Attributes struct {
	Name        string `json:"gema.service"`
	Proto       string `json:"gema.proto"`
	Port        string `json:"gema.port"`
	Auth        string `json:"gema.auth"`
	AccessLevel string `json:"gema.access_level"`
	Domain      string `json:"gema.domain"`
	SubDomain   string `json:"gema.subdomain"`
	Path        string `json:"gema.path"`
	CORS        string `json:"gema.cors"`
}

func DefaultEvent() *Event {
	return &Event{
		Status: "create",
		Actor: Actor{
			Attributes: Attributes{
				Name:        "",
				Proto:       "http",
				Port:        "",
				Auth:        "0",
				AccessLevel: "0",
				Domain:      "",
				SubDomain:   "",
				Path:        "/",
				CORS:		 "no",
			},
		},
	}
}

type Actor struct {
	Attributes Attributes
}

type Event struct {
	Status string `json:"status"`
	Actor  Actor
}

type Handler struct {
	Database *redis.Client
}

type ConfigRefreshEvent struct {
	Id string `json:"id"`
	Service string `json:"service"`
}

func New() *Handler {
	return &Handler{
		Database: database.New(),
	}
}

func (h *Handler) HandleEvent(ev *Event) {
	if ev.Status != "create" {
		return
	}

	log.Printf("Received event: %s of service %s", ev.Status, ev.Actor.Attributes.Name)

	evService := ev.Actor.Attributes

	route := getRoute(evService)
	key := fmt.Sprintf("service:%s", route)

	evService.Domain = route

	j := string(utils.ToJSON(evService))

	err := h.Database.Set(key, j, 0).Err()
	utils.Handle(err)

	log.Printf("Created route: %s for service: %s", route, evService.Name)

	refreshEvent := &ConfigRefreshEvent{
		Id: fmt.Sprint(time.Now().Unix()),
		Service: key,
	}

	// Push the config refresh event to redis, set expiration to give time for loading in proxy replicas.
	h.Database.LPush("service:events", utils.ToJSON(refreshEvent))
	h.Database.Expire("service:events", 30 * time.Second)
}

func getRoute(svc Attributes) string {
	route := ""

	stg := ""
	if os.Getenv("ENVIRONMENT") == "stg" {
		stg = "stg."
	}

	sub := ""
	if svc.SubDomain != "" {
		sub = svc.SubDomain + "."
	}

	if svc.Domain != "" {
		if strings.Contains(svc.Domain, ".") {
			route = fmt.Sprintf("%s%s%s", sub, stg, svc.Domain)
		} else {
			route = fmt.Sprintf("%s%s%s.%s", sub, stg, svc.Domain, os.Getenv("HQ_DOMAIN"))
		}
	} else {
		route = fmt.Sprintf("%s%s", sub, os.Getenv("HQ_DOMAIN"))
	}

	if os.Getenv("ENVIRONMENT") == "dev" {
		hq := "hq."
		if svc.Domain != "" {
			hq = svc.Name + "."
		}
		route = fmt.Sprintf("%s%s%s", sub, hq, "localhost")
	}

	return route
}
