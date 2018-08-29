package main

import (
	"github.com/kataras/iris"
	"gema/server/server"
	"gema/server/views"
	"gema/server/security"
	"gema/server/models"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/middleware/logger"
)

type Health struct {
	Status string
}

var gema *server.Server

func main() {
	app := iris.Default()

	app.Logger().Info("Starting GEMA on 8080.")

	app.Use(recover.New())

	l := logger.New(logger.Config{
		Status:   true,
		IP:       true,
		Method:   true,
		Path:     true,
		Query:    false,
		Columns:  false,
		LogFunc:  nil,
		Skippers: []logger.SkipperFunc{
			func(ctx iris.Context) bool {
				return ctx.Path() == "/gema/health"
			},
		},
	})

	app.Use(l)

	gema = server.New(app)

	iris.RegisterOnInterrupt(func() {
		gema.Dispose()
	})

	app.RegisterView(iris.HTML("./templates/landing", ".html").Layout("layout.html"))

	app.StaticWeb("/static", "./static")

	app.Get("/gema/health", healthHandler)

	app.Post("/gema/login", loginPostHandler)

	app.Get("/gema/setup", setupGetHandler)
	app.Post("/gema/setup", setupPostHandler)

	app.Handle("ALL", "*", proxyHandler)

	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
}

func healthHandler(ctx iris.Context) {
	ctx.JSON(&Health{
		Status: "OK",
	})
}

func loginPostHandler(ctx iris.Context) {
	email := ctx.PostValue("email")
	password := ctx.PostValue("password")

	user := models.GetUser(gema.Postgres, email)

	if user == nil {
		ctx.ViewData("Invalid", true)
		views.LoginPage(ctx)
	}

	login := security.ComparePasswords(user.Hash, password)

	if login {
		ctx.Application().Logger().Infof("%s logged in.", email)

		s := gema.Session.Start(ctx)
		s.Set("authorized", true)
		ctx.Redirect("/")
	}

	ctx.ViewData("Invalid", true)
	views.LoginPage(ctx)
}

func setupGetHandler(ctx iris.Context) {
	ctx.Application().Logger().Info("Setting up admin account.")
	views.SetupPage(ctx)
}

func setupPostHandler(ctx iris.Context) {
	email := ctx.PostValue("email")
	name := ctx.PostValue("name")
	password := ctx.PostValue("password")
	hash := security.GetHash(password)

	gema.Postgres.Insert(&models.User{
		Email: email,
		Name: name,
		Hash: hash,
	})

	ctx.Application().Logger().Info("Admin account ready.")

	ctx.Redirect("/")
}

func proxyHandler(ctx iris.Context) {
	if gema.IsFirstLogin() {
		ctx.Redirect("/gema/setup")
	}

	s := gema.Session.Start(ctx)

	if s.GetBooleanDefault("authorized", false) {
		ctx.Application().Logger().Info(ctx.Host())

		//target, _ := url.Parse("http://kibana:5601")
		//proxy := host.ProxyHandler(target)
		//proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	}

	views.LoginPage(ctx)
}
