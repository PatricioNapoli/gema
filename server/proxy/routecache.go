package proxy

import (
	"fmt"
	"gema/server/services"
	"gema/server/utils"
	"github.com/go-redis/redis"
	"time"
)

type RouteCache struct {
	services *services.Services

	cache map[string]string
	processed map[string]bool
}

func NewRouteCache(services *services.Services) *RouteCache {
	rc := &RouteCache{
		services:services,
		cache:make(map[string]string),
		processed:make(map[string]bool),
	}

	// Check for a config refresh event in Redis sent by the Agent
	go utils.DoEvery(5 * time.Second, rc.configChecker)

	// Prevents a huge array if application runs for a long time
	go utils.DoEvery(60 * time.Second, rc.clearProcessed)

	return rc
}

type ConfigRefreshEvent struct {
	Id string `json:"id"`
	Service string `json:"service"`
}

func (rc *RouteCache) configChecker() {
	events, err := rc.services.Store.LRange("service:events", 0, -1).Result()
	if err == redis.Nil {
		return
	}

	for _, ev := range events {
		event := &ConfigRefreshEvent{}
		utils.FromJSON([]byte(ev), event)

		if _, exists := rc.processed[event.Id]; !exists {
			svc, _ := rc.services.Store.Get(event.Service).Result()
			rc.cache[event.Service] = svc
			rc.processed[event.Id] = true
		}
	}
}

func (rc *RouteCache) clearProcessed() {
	rc.processed = make(map[string]bool)
}

func (rc *RouteCache) GetRouteConfig(host string) string {
	key := fmt.Sprintf("service:%s", host)

	if val, exists := rc.cache[key]; exists {
		return val
	}

	svc, err := rc.services.Store.Get(key).Result()
	if err == redis.Nil {
		return ""
	}

	rc.cache[key] = svc

	return svc
}