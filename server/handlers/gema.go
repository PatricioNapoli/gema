package handlers

import (
	"gema/server/models"
	"gema/server/security"
	"gema/server/views"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"gema/server/database"
	"github.com/kataras/iris/sessions"
)

type GEMA struct {
	NoSQL    *redis.Client
	Database *database.Database
	Session  *sessions.Sessions
}

type health struct {
	Status string
}

func (s *GEMA) Health(ctx iris.Context) {
	ctx.JSON(&health{
		Status: "OK",
	})
}

func (s *GEMA) LoginPost(ctx iris.Context) {
	email := ctx.PostValue("email")
	password := ctx.PostValue("password")

	user := models.FetchUserByEmail(s.Database.SQL, email)

	if user == nil {
		ctx.ViewData("Invalid", true)
		views.LoginPage(ctx)
		return
	}

	login := security.ComparePasswords(user.Hash, password)

	if login {
		ctx.Application().Logger().Infof("%s logged in.", email)

		s := s.Session.Start(ctx)
		s.Set("authorized", true)
		ctx.Redirect("/")
		return
	}

	ctx.ViewData("Invalid", true)
	views.LoginPage(ctx)
}

func (s *GEMA) SetupGet(ctx iris.Context) {
	if !s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/login")
		return
	}

	ctx.Application().Logger().Info("Setting up admin account.")
	views.SetupPage(ctx)
}

func (s *GEMA) SetupPost(ctx iris.Context) {
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

