package handlers

import (
	"gema/server/models"
	"gema/server/security"
	"gema/server/views"

	"fmt"
	"gema/server/services"

	"github.com/kataras/iris"
)

type Handlers struct {
	Services  *services.Services
	Dashboard *Dashboard
}

func New(services *services.Services) *Handlers {
	return &Handlers{
		Services: services,
		Dashboard: &Dashboard{
			Services: services,
		},
	}
}

func (s *Handlers) LoginPost(ctx iris.Context) {
	email := ctx.PostValue("email")
	password := ctx.PostValue("password")

	user := models.FetchUserByEmail(s.Services.Database.SQL, email)

	if user == nil {
		ctx.ViewData("Invalid", true)
		views.LoginPage(ctx)
		return
	}

	login := security.ComparePasswords(user.Hash, password)

	if login {
		ctx.Application().Logger().Infof("%s logged in.", email)

		s := s.Services.Session.Start(ctx)
		s.Set("authorized", true)

		if ctx.URLParamExists("next") {
			ctx.Redirect(fmt.Sprintf("https://%s/", ctx.URLParam("next")))
			return
		}

		ctx.Redirect("/")
		return
	}

	ctx.ViewData("Invalid", true)
	views.LoginPage(ctx)
}

func (s *Handlers) SetupGet(ctx iris.Context) {
	if !s.Services.Database.IsFirstLogin() {
		ctx.Redirect("/")
		return
	}

	ctx.Application().Logger().Info("Setting up admin account.")
	views.SetupPage(ctx)
}

func (s *Handlers) SetupPost(ctx iris.Context) {
	if !s.Services.Database.IsFirstLogin() {
		ctx.Redirect("/")
		return
	}

	email := ctx.PostValue("email")
	name := ctx.PostValue("name")
	password := ctx.PostValue("password")
	hash := security.GetHash(password)

	s.Services.Database.SQL.Insert(&models.User{
		Email: email,
		Name:  name,
		Hash:  hash,
	})

	ctx.Application().Logger().Info("Admin account ready.")

	ctx.Redirect("/")
}
