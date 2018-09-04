package main

import (
	"github.com/kataras/iris"
	"gema/server/server"
	"gema/server/views"
	"gema/server/security"
	"gema/server/models"
	"github.com/kataras/iris/middleware/logger"
	ravenIris "github.com/iris-contrib/middleware/raven"
	"regexp"
	"github.com/getsentry/raven-go"
	"github.com/elastic/apm-agent-go"
	"fmt"
)

type Health struct {
	Status string
}

var gema *server.Server

func main() {
	raven.SetDSN("http://c738fc18de3e4351b58ef0130899ab1e:479d921d673d4c9dbc1f78696d85eaba@sentry:9000/1")

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

	gema = server.New(app)

	iris.RegisterOnInterrupt(func() {
		gema.Dispose()
	})

	app.RegisterView(iris.HTML("./templates/landing", ".html").Layout("layout.html"))
	app.StaticWeb("/static", "./static")

	gemaRoute := app.Party("/gema")
	gemaRoute.Get("/health", healthHandler)
	gemaRoute.Post("/login", loginPostHandler)
	gemaRoute.Get("/setup", setupGetHandler)
	gemaRoute.Post("/setup", setupPostHandler)

	app.Handle("ALL", "*", proxyHandler)

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
	if !gema.IsFirstLogin() {
		ctx.Redirect("/gema/login")
	}

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
	tx := elasticapm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s", ctx.Method(), ctx.Path()), "proxy")
	tx.End()

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
