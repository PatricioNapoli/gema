package utils

import (
	"gema/server/views"
	"regexp"

	ravenIris "github.com/iris-contrib/middleware/raven"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

var StaticFilesRegex = regexp.MustCompile(`.+\..{2,5}.*$`)

func RegisterRecovery(app *iris.Application) {
	app.Logger().Info("Setting up Raven recovery handler.")

	app.Use(ravenIris.RecoveryHandler)
}

func RegisterLogger(app *iris.Application) {
	app.Logger().Info("Setting up logger.")

	app.Use(logger.New(logger.Config{
		Status:   true,
		IP:       true,
		Method:   true,
		Path:     true,
		Query:    true,
		Columns:  false,
		LogFunc:  nil,
		Skippers: []logger.SkipperFunc{logSkipper},
	}))
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

func MatchStaticFiles(path string) bool {
	return StaticFilesRegex.MatchString(path)
}

func logSkipper(ctx iris.Context) bool {
	if ctx.Path() == "/health" {
		return true
	}

	// Ignore static files
	matched := MatchStaticFiles(ctx.Path())

	return matched
}
