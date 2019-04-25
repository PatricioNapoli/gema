package utils

import (
	"gema/server/views"

	ravenIris "github.com/iris-contrib/middleware/raven"
	"github.com/kataras/iris"
)

func RegisterRecovery(app *iris.Application) {
	app.Logger().Info("Setting up Raven recovery handler.")

	app.Use(ravenIris.RecoveryHandler)
}

func RegisterViews(app *iris.Application) {
	app.Logger().Info("Setting up views.")

	app.RegisterView(iris.HTML("./templates", ".html").Layout("landing/landing_layout.html"))
}

func RegisterErrorHandlers(app *iris.Application) {
	app.Logger().Info("Setting up error handlers.")

	app.OnErrorCode(iris.StatusBadGateway, func(ctx iris.Context) {
		views.BadGateway(ctx)
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		views.InternalError(ctx)
	})

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		views.NotFound(ctx)
	})

	app.OnErrorCode(iris.StatusForbidden, func(ctx iris.Context) {
		views.Forbidden(ctx)
	})
}
