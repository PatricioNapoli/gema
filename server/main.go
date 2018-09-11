package main

import (
	"gema/server/server"
	"gema/server/services"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/kataras/iris"
	"gema/server/proxy"
)

func main() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))

	svcs := services.New()

	proxyApp := iris.New()
	prox := proxy.New(proxyApp, svcs)
	go prox.Start(":80")

	gemaApp := iris.New()
	gema := server.New(gemaApp, svcs)
	gema.Start(":81")

	iris.RegisterOnInterrupt(func() {
		svcs.Dispose()
	})
}
