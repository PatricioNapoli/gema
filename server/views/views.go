package views

import "github.com/kataras/iris"

func LoginPage(ctx iris.Context) {
	ctx.View("login.html")
}

func SetupPage(ctx iris.Context) {
	ctx.View("setup.html")
}