package main

import (
	"github.com/kataras/iris"
	//"gema/server/handlers"
	"gema/server/server"
	//"net/url"
	//"github.com/kataras/iris/core/host"
	"gema/server/handlers"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	app := iris.Default()
	app.Logger().Info("Starting GEMA on 8080.")

	gema := server.New(app)

	iris.RegisterOnInterrupt(func() {
		gema.Dispose()
	})

	app.RegisterView(iris.HTML("./views", ".html").Layout("layout.html"))

	app.StaticWeb("/static", "./static")

	app.Post("/login", func(ctx iris.Context) {
		//email := ctx.PostValue("email")
		password := ctx.PostValue("password")
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		app.Logger().Info(string(hash))
		app.Logger().Info(err)
		// check login correct, set session, then redirect to /

		ctx.Redirect("/")
	})

	app.Handle("ALL", "/*", func(ctx iris.Context) {
		s := gema.Session.Start(ctx)
		//s.Set("authorized", true)

		if s.GetBooleanDefault("authorized", false) {
			//target, _ := url.Parse("http://kibana:5601")
			//proxy := host.ProxyHandler(target)
			//proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		}

		//ctx.ViewData("Name", s.GetString("name"))
		handlers.LoginPage(ctx)
	})

	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.Configuration{ // default configuration:
		DisableStartupLog:                 true,
	}))
}
