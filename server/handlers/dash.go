package handlers

import (
	"fmt"
	"gema/server/services"
	"gema/server/views"
	"github.com/kataras/iris"
	"os"
)

type Dashboard struct {
	Services *services.Services
}

func (s *Dashboard) HQ(ctx iris.Context) {
	if !s.requiresLogin(ctx) {
		views.HQ(ctx)
	}
}

func (s *Dashboard) Purge(ctx iris.Context) {
	if !s.requiresLogin(ctx) {
		os.RemoveAll("/cache")
		ctx.Redirect("/")
	}
}

func (s *Dashboard) requiresLogin(ctx iris.Context) bool {
	session := s.Services.Session.Start(ctx)

	if session.GetBooleanDefault("authenticated", false) {
		return false
	}

	if ctx.URLParamExists("next") {
		ctx.ViewData("next", fmt.Sprintf("?next=%s", ctx.URLParam("next")))
	}

	views.LoginPage(ctx)
	return true
}

// TODO: Clear another user's session?