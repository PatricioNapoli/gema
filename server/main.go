package main

import (
	"gema/server/server"
	"gema/server/views"
	"os"
	"regexp"

	"github.com/getsentry/raven-go"
	ravenIris "github.com/iris-contrib/middleware/raven"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func main() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))

	app := iris.New()

	app.Use(ravenIris.RecoveryHandler)

	l := logger.New(logger.Config{
		Status:   true,
		IP:       true,
		Method:   true,
		Path:     true,
		Query:    true,
		Columns:  false,
		LogFunc:  nil,
		Skippers: []logger.SkipperFunc{logSkipper},
	})

	app.Use(l)

	gema := server.New(app)

	iris.RegisterOnInterrupt(func() {
		gema.Dispose()
	})

	app.RegisterView(iris.HTML("./templates", ".html").Layout("landing/landing_layout.html"))
	app.StaticWeb("/static/gema", "./static/gema")

	app.OnErrorCode(iris.StatusBadGateway, func(ctx iris.Context) {
		views.InternalError(ctx)
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		views.InternalError(ctx)
	})

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		views.NotFound(ctx)
	})

	app.OnErrorCode(iris.StatusForbidden, func(ctx iris.Context) {
		views.InternalError(ctx)
	})

	gemaRoute := app.Party("/gema")
	gemaRoute.Get("/health", gema.Handlers.Health)
	gemaRoute.Post("/login", gema.Handlers.LoginPost)
	gemaRoute.Get("/setup", gema.Handlers.SetupGet)
	gemaRoute.Post("/setup", gema.Handlers.SetupPost)

	app.Handle("ALL", "*", gema.Handlers.Proxy)

	app.Logger().Info("Starting GEMA server.")
	app.Run(iris.Addr(":80"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
		DisablePathCorrection: true,
	}))
}

func logSkipper(ctx iris.Context) bool {
	if ctx.Path() == "/gema/health" {
		return true
	}

	matched, err := regexp.MatchString(`.+\..{2,5}$`, ctx.Path())

	if err != nil {
		panic(err)
	}

	return matched
}
