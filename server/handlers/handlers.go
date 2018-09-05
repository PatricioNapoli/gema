package handlers

import (
	"github.com/kataras/iris"
	"gema/server/models"
	"gema/server/views"
	"gema/server/security"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions"
	"gema/server/database"
	"github.com/elastic/apm-agent-go"
	"fmt"
)

type Handlers struct {
	NoSQL   *redis.Database
	Database     *database.Database
	Session *sessions.Sessions
}

func (s *Handlers) Dispose() {
	s.NoSQL.Close()
	s.Database.Dispose()
}

type Health struct {
	Status string
}

func (s *Handlers) Health(ctx iris.Context) {
	ctx.JSON(&Health{
		Status: "OK",
	})
}

func (s *Handlers) Proxy(ctx iris.Context) {
	if s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/setup")
	}

	session := s.Session.Start(ctx)

	if session.GetBooleanDefault("authorized", false) {
		ctx.Application().Logger().Info(ctx.Host())

		tx := elasticapm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s", ctx.Method(), ctx.Path()), "proxy")
		defer tx.End()

		//target, _ := url.Parse("http://kibana:5601")
		//proxy := host.ProxyHandler(target)
		//proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	}

	views.LoginPage(ctx)
}

func (s *Handlers) LoginPost(ctx iris.Context) {
	email := ctx.PostValue("email")
	password := ctx.PostValue("password")

	user := &models.User{}
	user.FetchUserByEmail(s.Database.SQL, email)

	if user == nil {
		ctx.ViewData("Invalid", true)
		views.LoginPage(ctx)
	}

	login := security.ComparePasswords(user.Hash, password)

	if login {
		ctx.Application().Logger().Infof("%s logged in.", email)

		s := s.Session.Start(ctx)
		s.Set("authorized", true)
		ctx.Redirect("/")
	}

	ctx.ViewData("Invalid", true)
	views.LoginPage(ctx)
}

func (s *Handlers) SetupGet(ctx iris.Context) {
	if !s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/login")
	}

	ctx.Application().Logger().Info("Setting up admin account.")
	views.SetupPage(ctx)
}

func (s *Handlers) SetupPost(ctx iris.Context) {
	email := ctx.PostValue("email")
	name := ctx.PostValue("name")
	password := ctx.PostValue("password")
	hash := security.GetHash(password)

	s.Database.SQL.Insert(&models.User{
		Email: email,
		Name: name,
		Hash: hash,
	})

	ctx.Application().Logger().Info("Admin account ready.")

	ctx.Redirect("/")
}