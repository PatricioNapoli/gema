package services

import (
	"gema/server/database"
	"os"
	"time"

	"go.elastic.co/apm"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/sessions"
	irisRedis "github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
)

type Services struct {
	Database *database.Database
	NoSQL    *redis.Client
	Session  *sessions.Sessions
	Tracing *apm.Tracer
}

// Creates a new Services object containing all the services needed for operating with environment.
func New() *Services {
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
		Expires:      time.Duration(24 * 7) * time.Hour,
		AllowReclaim: true,
	})

	s.UseDatabase(r)

	if d.Migrate() {
		// DB was droppped, destroy sessions
		s.DestroyAll()
	}

	return &Services{
		Database: d,
		NoSQL:    rc,
		Session:  s,
		Tracing: apm.DefaultTracer,
	}
}

func (s *Services) Dispose() {
	s.Database.Dispose()
	s.NoSQL.Close()
}
