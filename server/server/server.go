package server

import (
	"github.com/go-pg/pg"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"time"
	"os"
	"github.com/go-pg/pg/orm"
	"gema/server/models"
)

type Server struct {
	App *iris.Application
	Redis *redis.Database
	Postgres *pg.DB
	Session *sessions.Sessions
}

func New(app *iris.Application) *Server {
	app.Logger().Info("Initializing cache and main databases.")

	r := redis.New(service.Config{
		Network:     "tcp",
		Addr:        "redis:6379",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: time.Duration(5) * time.Minute,
		Prefix:      "gema:"})

	p := pg.Connect(&pg.Options{
		Addr: "postgres:5432",
		User: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_USER"),
	})

	migrate(p)

	s := sessions.New(sessions.Config{
		Cookie:       "sessionscookieid",
		Expires:      0,
		AllowReclaim: true,
	})

	s.UseDatabase(r)

	app.Logger().Info("Databases initialized.")

	return &Server{
		App: app,
		Redis: r,
		Postgres : p,
		Session: s,
	}
}

func (s *Server) IsFirstLogin() bool {
	count, err := s.Postgres.Model((*models.User)(nil)).Count()
	if err != nil {
		panic(err)
	}
	return count == 0
}

func (s Server) Dispose() {
	s.Redis.Close()
	s.Postgres.Close()
}

func migrate(db *pg.DB) error {
	for _, model := range []interface{}{(*models.User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}