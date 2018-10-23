package handlers

import (
	"github.com/kataras/iris"
	"gema/server/services"
	"gema/server/views"
	"fmt"
)

type Dashboard struct {
	Services *services.Services
}

func (s *Dashboard) HQ(ctx iris.Context) {
	s.requireLogin(ctx)

	views.HQ(ctx)
}

func (s *Dashboard) requireLogin(ctx iris.Context) {
	session := s.Services.Session.Start(ctx)

	if session.GetBooleanDefault("authenticated", false) {
		return
	}

	if ctx.URLParamExists("next") {
		ctx.ViewData("next", fmt.Sprintf("?next=%s", ctx.URLParam("next")))
	}

	views.LoginPage(ctx)
}

// TODO: Logout
// TODO: Clear another user's session?