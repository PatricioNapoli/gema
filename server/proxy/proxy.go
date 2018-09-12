package proxy

import (
	"crypto/tls"
	"fmt"
	"gema/server/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/netutil"
	"gema/server/services"
)

var (
	WebDAVMethods = [...]string{
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
	App *iris.Application
	Services *services.Services
}

// Creates a new Proxy which receives the requests and proxies them based on the config found in redis.
func New(app *iris.Application, services *services.Services) *Proxy {
	app.Logger().Info("Setting up GEMA reverse proxy.")

	utils.RegisterRecovery(app)
	utils.RegisterLogger(app)
	utils.RegisterViews(app)
	utils.RegisterErrorHandlers(app)

	proxy := &Proxy{
		App: app,
		Services: services,
	}

	proxy.setupRoutes()

	return proxy
}

func (s *Proxy) Start(addr string) {
	s.App.Logger().Info("Starting GEMA proxy.")

	s.App.Run(iris.Addr(addr), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:     true,
		DisablePathCorrection: true,
		DisableVersionChecker: true,
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
	Name      string `json:"name"`
	Proto     string `json:"proto"`
	Port      string `json:"port"`
	Auth      string `json:"auth"`
	Domain    string `json:"domain"`
	SubDomain string `json:"subdomain"`
	Path      string `json:"path"`
}

func (s *Proxy) proxy(ctx iris.Context) {
	session := s.Services.Session.Start(ctx)

	if ctx.Host() == os.Getenv("HQ_DOMAIN") {
		target, _ := url.Parse("http://localhost:81/")
		proxy := proxyHandler(target, ctx.Host(), ctx.GetHeader("X-Real-IP"))
		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		return
	}

	svc, err := s.Services.NoSQL.Get(fmt.Sprintf("service:%s", ctx.Host())).Result()
	if err == redis.Nil {
		ctx.NotFound()
		return
	}

	serv := &service{}
	utils.FromJSON([]byte(svc), &serv)

	if session.GetBooleanDefault("authorized", false) || serv.Auth == "0" {
		tx := s.Services.Tracing.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), "proxy")
		defer tx.End()

		port := ":80"
		if serv.Port != "" {
			port = fmt.Sprintf(":%s", serv.Port)
		}

		route := fmt.Sprintf("%s://%s%s%s", serv.Proto, serv.Name, port, serv.Path)

		target, _ := url.Parse(route)
		proxy := proxyHandler(target, ctx.Host(), ctx.GetHeader("X-Real-IP"))
		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		return
	}

	ctx.Redirect(fmt.Sprintf("https://%s/?next=%s", os.Getenv("HQ_DOMAIN"), ctx.Host()))
}


// This proxy handler has almost the same functionalty as the original net/http proxy handler.
// What's different, is that it has an added feature which receives a custom host to forward.
// Also, sets the common reverse proxy headers, X-Forwarded-Host and X-Real-IP.
func proxyHandler(target *url.URL, originalHost string, realIp string) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = originalHost
		req.Header.Set("X-Forwarded-Host", originalHost)
		req.Header.Set("X-Real-IP", realIp)
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	p := &httputil.ReverseProxy{Director: director}

	if netutil.IsLoopbackHost(target.Host) {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		p.Transport = transport
	}

	return p
}

// Added here because it was private in net/http.
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
