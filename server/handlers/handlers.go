package handlers

import (
	"github.com/go-redis/redis"
	"gema/server/database"
	"github.com/kataras/iris/sessions"
	"fmt"
	"github.com/kataras/iris"
	"os"
	"gema/server/views"
	"github.com/elastic/apm-agent-go"
	"net/url"
	"net/http/httputil"
	"net/http"
	"github.com/kataras/iris/core/netutil"
	"crypto/tls"
	"strings"
	"gema/server/utils"
)

type Handlers struct {
	NoSQL    *redis.Client
	Database *database.Database
	Session  *sessions.Sessions
	GEMA *GEMA
	Dashboard *Dashboard
}

func New(nosql *redis.Client, database *database.Database, session *sessions.Sessions) *Handlers {
	return &Handlers{
		NoSQL: nosql,
		Database: database,
		Session: session,
		GEMA: &GEMA{
			NoSQL: nosql,
			Database: database,
			Session: session,
		},
		Dashboard: &Dashboard{
			NoSQL: nosql,
			Database: database,
			Session: session,
		},
	}
}

func (s *Handlers) Dispose() {
	s.NoSQL.Close()
	s.Database.Dispose()
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

func (s *Handlers) Proxy(ctx iris.Context) {
	if s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/setup")
	}

	session := s.Session.Start(ctx)

	ctx.Application().Logger().Info("Redirecting to dash: " + ctx.Host())

	if ctx.Host() == os.Getenv("HQ_DOMAIN") {
		if session.GetBooleanDefault("authorized", false) {
			ctx.Application().Logger().Info("Redirecting to dash: " + ctx.Host())
			ctx.Redirect("/dash/view	")
			return
		}

		views.LoginPage(ctx)
		return
	}

	svc, err := s.NoSQL.Get(fmt.Sprintf("service:%s", ctx.Host())).Result()
	if err == redis.Nil {
		ctx.NotFound()
		return
	}

	serv := &service{}
	utils.FromJSON([]byte(svc), &serv)

	if session.GetBooleanDefault("authorized", false) || serv.Auth == "0" {
		tx := elasticapm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), "proxy")
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

	views.LoginPage(ctx)
}

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