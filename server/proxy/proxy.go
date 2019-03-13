package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gema/server/utils"
	"go.elastic.co/apm"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gema/server/services"

	"github.com/kataras/iris"
)

var (
	WebDAVMethods = [...]string{
		"REPORT",
		"COPY",
		"LOCK",
		"MKCOL",
		"MOVE",
		"PROPFIND",
		"PROPPATCH",
		"UNLOCK",
	}
)

// Wrapper with Application pointer and services.
type Proxy struct {
	App      *iris.Application
	Services *services.Services

	routeCache *RouteCache
}

// Creates a new Proxy which receives the requests and proxies them based on the config found in redis.
func New(app *iris.Application, services *services.Services) *Proxy {
	app.Logger().Info("Setting up GEMA reverse proxy.")

	utils.RegisterRecovery(app)
	utils.RegisterLogger(app)
	utils.RegisterViews(app)
	utils.RegisterErrorHandlers(app)

	proxy := &Proxy{
		App:      app,
		Services: services,

		routeCache: NewRouteCache(services),
	}

	proxy.setupRoutes()

	return proxy
}

func (s *Proxy) Start(addr string) {
	s.App.Logger().Info("Starting GEMA proxy.")

	s.App.Run(iris.Addr(addr), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:     true,
		DisablePathCorrection: true,
	}))
}

func (s *Proxy) setupRoutes() {
	s.App.Logger().Info("Setting up reverse proxy handlers.")

	s.App.Handle("ALL", "*", s.proxy)

	// WebDAV hack
	for _, method := range WebDAVMethods {
		s.App.Handle(method, "*", s.proxy)
	}
}

type service struct {
	Name        string `json:"gema.service"`
	Proto       string `json:"gema.proto"`
	Port        string `json:"gema.port"`
	Auth        string `json:"gema.auth"`
	AccessLevel string `json:"gema.access_level"`
	Domain      string `json:"gema.domain"`
	SubDomain   string `json:"gema.subdomain"`
	Path        string `json:"gema.path"`
}

func (s *Proxy) proxy(ctx iris.Context) {
	session := s.Services.Session.Start(ctx)

	if ctx.Host() == os.Getenv("HQ_DOMAIN") {
		if !strings.Contains(ctx.Path(), "websocket") {
			tx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), "hq")
			defer tx.End()
		}

		target, _ := url.Parse("http://localhost:81/")
		proxy := NewHTTPProxy(target, ctx.Host(), ctx.GetHeader("X-Real-IP"))
		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		return
	}

	svc := s.routeCache.GetRouteConfig(ctx.Host())
	if svc == "" {
		ctx.NotFound()
		return
	}

	serv := &service{}
	utils.FromJSON([]byte(svc), &serv)

	if !strings.Contains(ctx.Path(), "websocket") {
		tx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), serv.Name)
		defer tx.End()
	}

	if serv.Auth == "0" || session.GetBooleanDefault("authenticated", false)  {
		port := ":80"
		if serv.Port != "" {
			port = fmt.Sprintf(":%s", serv.Port)
		}

		proto := serv.Proto
		if ctx.Request().Header.Get("Connection") == "upgrade" {
			proto = "ws"
		}

		route := fmt.Sprintf("%s://%s%s%s", proto, serv.Name, port, serv.Path)

		target, _ := url.Parse(route)

		// Handle WebSocket
		if ctx.Request().Header.Get("Connection") == "upgrade" {
			wsProxy := NewWebSocketProxy(target, ctx.Host(), ctx.GetHeader("X-Real-IP"))
			wsProxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			return
		}

		proxy := NewHTTPProxy(target, ctx.Host(), ctx.GetHeader("X-Real-IP"))

		if utils.MatchStaticFiles(ctx.Path()) {
			proxy.ModifyResponse = interceptStaticFile
		}

		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())

		return
	}

	ctx.Redirect(fmt.Sprintf("https://%s/?next=%s", os.Getenv("HQ_DOMAIN"), ctx.Host()))
}

func interceptStaticFile(resp *http.Response) (err error) {
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
		if err != nil {
			return err
		}
	} else {
		reader = resp.Body
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return  err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	fileDir := fmt.Sprintf("/static/%s/%s", resp.Request.Host, resp.Request.RequestURI[1:])
	err = os.MkdirAll(filepath.Dir(fileDir), 0755)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileDir, b, 0644)
	if err != nil {
		return err
	}

	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	return nil
}
