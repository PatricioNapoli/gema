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
	ravenIris "github.com/iris-contrib/middleware/raven"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"github.com/getsentry/raven-go"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/context"
	"regexp"
	"gema/server/views"
)

var (
	WebDAVMethods = [...]string{
		"COPY",
		"LOCK",
		"MKCOL",
		"MOVE",
		"PROPFIND",
		"PROPPATCH",
		"UNLOCK",
	}
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
		Handlers: handlers.New(rc, d, s),
	}
}

func (s Server) Dispose() {
	s.Handlers.Dispose()
}

func (s Server) Setup() {
	s.App.Logger().Info("Setting up GEMA server.")

	raven.SetDSN(os.Getenv("SENTRY_DSN"))

	s.App.Logger().Info("Setting up recovery handler and logger.")
	s.App.Use(ravenIris.RecoveryHandler)
	s.App.Use(getLogger())

	iris.RegisterOnInterrupt(func() {
		s.Dispose()
	})

	s.registerStatic()
	s.registerErrorHandlers()
	s.setupReverseProxy()
	s.setupRoutes()
}

func (s Server) Start() {
	s.App.Logger().Info("Starting GEMA server.")

	s.App.Run(iris.Addr(":80"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:     true,
		DisablePathCorrection: true,
	}))
}

func (s Server) registerStatic() {
	s.App.Logger().Info("Setting up static routes and views.")

	s.App.RegisterView(iris.HTML("./templates", ".html").Layout("landing/landing_layout.html"))
	s.App.StaticWeb("/static/gema", "./static/gema")
}

func (s Server) registerErrorHandlers() {
	s.App.Logger().Info("Setting up error handlers.")

	s.App.OnErrorCode(iris.StatusBadGateway, func(ctx iris.Context) {
		views.BadGateway(ctx)
	})

	s.App.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		views.InternalError(ctx)
	})

	s.App.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		views.NotFound(ctx)
	})

	s.App.OnErrorCode(iris.StatusForbidden, func(ctx iris.Context) {
		views.Forbidden(ctx)
	})
}

func (s Server) setupRoutes() {
	s.App.Logger().Info("Setting up routes.")

	gemaRoute := s.App.Party("/gema")
	gemaRoute.Get("/health", s.Handlers.GEMA.Health)
	gemaRoute.Post("/login", s.Handlers.GEMA.LoginPost)
	gemaRoute.Get("/setup", s.Handlers.GEMA.SetupGet)
	gemaRoute.Post("/setup", s.Handlers.GEMA.SetupPost)

	dashRoute := s.App.Party("/dash")
	dashRoute.Get("/view", s.Handlers.Dashboard.HQ)
}

func (s Server) setupReverseProxy() {
	s.App.Logger().Info("Setting up reverse proxy.")

	s.App.Handle("ALL", "*", s.Handlers.Proxy)

	// WebDAV hack
	for _, method := range WebDAVMethods {
		s.App.Handle(method, "*", s.Handlers.Proxy)
	}
}

func getLogger() context.Handler {
	return logger.New(logger.Config{
		Status:   true,
		IP:       true,
		Method:   true,
		Path:     true,
		Query:    true,
		Columns:  false,
		LogFunc:  nil,
		Skippers: []logger.SkipperFunc{logSkipper},
	})
}

func logSkipper(ctx iris.Context) bool {
	if ctx.Path() == "/gema/health" {
		return true
	}

	// Ignore static files
	matched, err := regexp.MatchString(`.+\..{2,5}.*$`, ctx.Path())

	if err != nil {
		panic(err)
	}

	return matched
}