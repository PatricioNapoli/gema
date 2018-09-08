package event

import (
	"fmt"
	"gema/agent/database"
	"gema/agent/utils"
	"log"
	"os"

	"github.com/go-redis/redis"
)

type Service struct {
	Name      string `json:"name"`
	Proto     string `json:"proto"`
	Port      string `json:"port"`
	Auth      string `json:"auth"`
	Domain    string `json:"domain"`
	SubDomain string `json:"subdomain"`
	Path      string `json:"path"`
}

type Attributes struct {
	Name      string `json:"gema.service"`
	Proto     string `json:"gema.proto"`
	Port      string `json:"gema.port"`
	Auth      string `json:"gema.auth"`
	Domain    string `json:"gema.domain"`
	SubDomain string `json:"gema.subdomain"`
	Path      string `json:"gema.path"`
}

func DefaultEvent() *Event {
	return &Event{
		Status: "create",
		Actor: Actor{
			Attributes: Attributes{
				Name:      "",
				Proto:     "http",
				Port:      "",
				Auth:      "0",
				Domain:    "",
				SubDomain: "",
				Path:      "/",
			},
		},
	}
}

func convertToService(attr Attributes) Service {
	return Service{
		Name:      attr.Name,
		Proto:     attr.Proto,
		Port:      attr.Port,
		Auth:      attr.Auth,
		Domain:    attr.Domain,
		SubDomain: attr.SubDomain,
		Path:      attr.Path,
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

	svc := convertToService(evService)

	j := string(utils.ToJSON(svc))

	route := getRoute(svc)

	err := h.Database.Set(fmt.Sprintf("service:%s", route), j, 0).Err()
	utils.Handle(err)

	log.Printf("Created route: %s for service: %s", route, svc.Name)
}

func getRoute(svc Service) string {
	route := ""

	stg := ""
	if os.Getenv("ENVIRONMENT") == "STAGING" {
		stg = "stg."
	}

	sub := ""
	if svc.SubDomain != "" {
		sub = svc.SubDomain + "."
	}

	if svc.Domain != "" {
		route = fmt.Sprintf("%s%s%s", sub, stg, svc.Domain)
	} else {
		route = fmt.Sprintf("%s%s", sub, os.Getenv("HQ_DOMAIN"))
	}

	if os.Getenv("ENVIRONMENT") == "DEVELOPMENT" {
		hq := "hq."
		if svc.Domain != "" {
			hq = svc.Domain + "."
		}
		route = fmt.Sprintf("%s%s%s", sub, hq, "localhost")
	}

	return route
}
