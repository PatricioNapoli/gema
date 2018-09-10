package handlers

import (
	"github.com/go-redis/redis"
	"gema/server/database"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris"
)

type Dashboard struct {
	NoSQL    *redis.Client
	Database *database.Database
	Session  *sessions.Sessions
}

func (s *Dashboard) HQ(ctx iris.Context) {
	ctx.WriteString("HQ")
}
