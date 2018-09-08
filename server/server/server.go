package server

import (
	"gema/server/database"
	"gema/server/handlers"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	irisRedis "github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
)

type Server struct {
	App      *iris.Application
	Handlers *handlers.Handlers
}

func New(app *iris.Application) *Server {
	app.Logger().Info("Initializing cache and main databases.")

	r := irisRedis.New(service.Config{
		Network:     "tcp",
		Addr:        "redis:6379",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: time.Duration(5) * time.Minute,
		Prefix:      "gema:"})

	rc := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	d := database.New(pg.Connect(&pg.Options{
		Addr:     "postgres:5432",
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_USER"),
	}))

	s := sessions.New(sessions.Config{
		Cookie:       "sessionscookieid",
		Expires:      0,
		AllowReclaim: true,
	})

	s.UseDatabase(r)

	app.Logger().Info("Databases initialized.")

	return &Server{
		App:      app,
		Handlers: &handlers.Handlers{Database: d, NoSQL: rc, Session: s},
	}
}

func (s Server) Dispose() {
	s.Handlers.Dispose()
}
