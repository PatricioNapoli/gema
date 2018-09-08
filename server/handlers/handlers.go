package handlers

import (
	"fmt"
	"net/url"

	"gema/server/database"
	"gema/server/models"
	"gema/server/security"
	"gema/server/utils"
	"gema/server/views"

	"github.com/elastic/apm-agent-go"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris/sessions"
)

type service struct {
	Name      string `json:"name"`
	Proto     string `json:"proto"`
	Port      string `json:"port"`
	Auth      string `json:"auth"`
	Domain    string `json:"domain"`
	SubDomain string `json:"subdomain"`
	Path      string `json:"path"`
}

type Handlers struct {
	NoSQL    *redis.Client
	Database *database.Database
	Session  *sessions.Sessions
}

func (s *Handlers) Dispose() {
	s.NoSQL.Close()
	s.Database.Dispose()
}

type health struct {
	Status string
}

func (s *Handlers) Health(ctx iris.Context) {
	ctx.JSON(&health{
		Status: "OK",
	})
}

func (s *Handlers) Proxy(ctx iris.Context) {
	if s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/setup")
	}

	tx := elasticapm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), "gateway")
	defer tx.End()

	session := s.Session.Start(ctx)

	ctx.Application().Logger().Info(ctx.GetHeader("Real-Host"))

	svc, err := s.NoSQL.Get(fmt.Sprintf("service:%s", ctx.Host())).Result()
	if err == redis.Nil {
		ctx.NotFound()
		return
	}

	serv := &service{}
	utils.FromJSON([]byte(svc), &serv)

	if session.GetBooleanDefault("authorized", false) || serv.Auth == "0" {
		ctx.Application().Logger().Info(ctx.Host())

		tx := elasticapm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), "proxy")
		defer tx.End()

		sub := ""
		if serv.SubDomain != "" {
			sub = fmt.Sprintf("%s.", serv.SubDomain)
		}
		port := ""
		if serv.Port != "" {
			port = fmt.Sprintf(":%s", serv.Port)
		}

		route := fmt.Sprintf("%s://%s%s%s%s", serv.Proto, sub, serv.Domain, port, serv.Path)

		ctx.Application().Logger().Info(fmt.Sprintf("Handling reverse proxy to %s", route))

		target, _ := url.Parse(route)
		proxy := host.ProxyHandler(target)
		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
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
		Name:  name,
		Hash:  hash,
	})

	ctx.Application().Logger().Info("Admin account ready.")

	ctx.Redirect("/")
}
