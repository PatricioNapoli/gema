package main

import (
	"github.com/kataras/iris"
	"gema/server/server"
	"github.com/kataras/iris/middleware/logger"
	ravenIris "github.com/iris-contrib/middleware/raven"
	"regexp"
	"github.com/getsentry/raven-go"
	"os"
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
		Query:    false,
		Columns:  false,
		LogFunc:  nil,
		Skippers: []logger.SkipperFunc{logSkipper},
	})

	app.Use(l)

	gema := server.New(app)

	iris.RegisterOnInterrupt(func() {
		gema.Dispose()
	})

	app.RegisterView(iris.HTML("./templates/landing", ".html").Layout("layout.html"))
	app.StaticWeb("/static", "./static")

	gemaRoute := app.Party("/gema")
	gemaRoute.Get("/health", gema.Handlers.Health)
	gemaRoute.Post("/login", gema.Handlers.LoginPost)
	gemaRoute.Get("/setup", gema.Handlers.SetupGet)
	gemaRoute.Post("/setup", gema.Handlers.SetupPost)

	app.Handle("ALL", "*", gema.Handlers.Proxy)

	app.Logger().Info("Starting GEMA server.")
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
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
