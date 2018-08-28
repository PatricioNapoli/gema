package handlers

import "github.com/kataras/iris"

func LoginPage(ctx iris.Context) {
	ctx.View("login.html")
}