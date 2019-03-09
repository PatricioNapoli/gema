package handlers

import (
	"gema/server/models"
	"gema/server/security"
	"gema/server/views"
	"time"

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

		sess := s.Services.Session.Start(ctx)
		sess.Set("authenticated", true)

		// Set last sign in
		user.LastSignIn = time.Now()
		s.Services.Database.SQL.Update(user)

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

func (s *Handlers) Logout(ctx iris.Context) {
	session := s.Services.Session.Start(ctx)
	session.Set("authenticated", false)

	ctx.Redirect("/")
}

func (s *Handlers) SetupGet(ctx iris.Context) {
	if !s.Services.Database.IsFirstUser() {
		ctx.Redirect("/")
		return
	}

	ctx.Application().Logger().Info("Setting up admin account.")
	views.SetupPage(ctx)
}

func (s *Handlers) SetupPost(ctx iris.Context) {
	if !s.Services.Database.IsFirstUser() {
		ctx.Redirect("/")
		return
	}

	email := ctx.PostValue("email")
	name := ctx.PostValue("name")
	password := ctx.PostValue("password")
	hash := security.GetHash(password)

	var uid int64
	s.Services.Database.SQL.Model(&models.User{
		Email: email,
		Name:  name,
		Hash:  hash,
	}).Returning("id", uid).Insert()

	s.Services.Database.SQL.Model(&models.Membership{
		UserId: uid,
		GroupId: models.FetchGroupByName(s.Services.Database.SQL, "master").Id,
	}).Insert()

	ctx.Application().Logger().Info("Admin account ready.")

	ctx.Redirect("/")
}
