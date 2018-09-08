package views

import "github.com/kataras/iris"

func LoginPage(ctx iris.Context) {
	ctx.View("login.html")
}

func SetupPage(ctx iris.Context) {
	ctx.View("setup.html")
}

func BadGateway(ctx iris.Context) {
	renderError(ctx, "502", "Service is not responding.")
}

func InternalError(ctx iris.Context) {
	renderError(ctx, "500", "Something went terribly wrong. We are working on the issue.")
}

func NotFound(ctx iris.Context) {
	renderError(ctx, "404", "We could not find what you were looking for.")
}

func Forbidden(ctx iris.Context) {
	renderError(ctx, "403", "You cannot perform this action.")
}

func renderError(ctx iris.Context, code string, description string) {
	ctx.ViewLayout("error/error_layout.html")
	ctx.ViewData("Code", code)
	ctx.ViewData("Description", description)
	ctx.View("error.html")
}
