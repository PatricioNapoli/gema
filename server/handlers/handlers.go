package handlers

import (
	"fmt"
	"net/url"
	"os"

	"gema/server/database"
	"gema/server/models"
	"gema/server/security"
	"gema/server/utils"
	"gema/server/views"

	"github.com/elastic/apm-agent-go"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"net/http/httputil"
	"net/http"
	"github.com/kataras/iris/core/netutil"
	"crypto/tls"
	"strings"
)

type service struct {
	Name      string `json:"name"`
	Proto     string `json:"proto"`
	Port      string `json:"port"`
	Auth      string `json:"auth"`
	Domain    string `json:"domain"`
	SubDomain string `json:"subdomain"`
	Path      string `json:"path"`
}

type Handlers struct {
	NoSQL    *redis.Client
	Database *database.Database
	Session  *sessions.Sessions
}

func (s *Handlers) Dispose() {
	s.NoSQL.Close()
	s.Database.Dispose()
}

type health struct {
	Status string
}

func (s *Handlers) Health(ctx iris.Context) {
	ctx.JSON(&health{
		Status: "OK",
	})
}

func (s *Handlers) Proxy(ctx iris.Context) {
	if s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/setup")
	}

	ctx.Application().Logger().Info(fmt.Sprintf("Handling reverse proxy from %s", ctx.Host()))

	tx := elasticapm.DefaultTracer.StartTransaction(fmt.Sprintf("%s %s%s", ctx.Method(), ctx.Host(), ctx.Path()), "gateway")
	defer tx.End()

	session := s.Session.Start(ctx)

	if ctx.Host() == os.Getenv("HQ_DOMAIN") {
		if session.GetBooleanDefault("authorized", false) {
			s.HQ(ctx)
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

		ctx.Application().Logger().Info(fmt.Sprintf("Handling reverse proxy to %s", route))

		target, _ := url.Parse(route)
		ctx.Application().Logger().Info(ctx.GetHeader("X-Real-IP"))
		proxy := proxyHandler(target, ctx.Host(), ctx.GetHeader("X-Real-IP"))
		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		return
	}

	views.LoginPage(ctx)
}

func (s *Handlers) LoginPost(ctx iris.Context) {
	email := ctx.PostValue("email")
	password := ctx.PostValue("password")

	user := models.FetchUserByEmail(s.Database.SQL, email)

	if user == nil {
		ctx.ViewData("Invalid", true)
		views.LoginPage(ctx)
		return
	}

	login := security.ComparePasswords(user.Hash, password)

	if login {
		ctx.Application().Logger().Infof("%s logged in.", email)

		s := s.Session.Start(ctx)
		s.Set("authorized", true)
		ctx.Redirect("/")
		return
	}

	ctx.ViewData("Invalid", true)
	views.LoginPage(ctx)
}

func (s *Handlers) SetupGet(ctx iris.Context) {
	if !s.Database.IsFirstLogin() {
		ctx.Redirect("/gema/login")
		return
	}

	ctx.Application().Logger().Info("Setting up admin account.")
	views.SetupPage(ctx)
}

func (s *Handlers) SetupPost(ctx iris.Context) {
	email := ctx.PostValue("email")
	name := ctx.PostValue("name")
	password := ctx.PostValue("password")
	hash := security.GetHash(password)

	s.Database.SQL.Insert(&models.User{
		Email: email,
		Name:  name,
		Hash:  hash,
	})

	ctx.Application().Logger().Info("Admin account ready.")

	ctx.Redirect("/")
}

func (s *Handlers) HQ(ctx iris.Context) {
	ctx.WriteString("HQ")
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