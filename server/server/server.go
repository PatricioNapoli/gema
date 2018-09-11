package server

import (
	"gema/server/handlers"
	"gema/server/services"
	"gema/server/utils"
	"github.com/kataras/iris"
)

type Server struct {
	App      *iris.Application
	Handlers *handlers.Handlers
}

func New(app *iris.Application, services *services.Services) *Server {
	app.Logger().Info("Setting up GEMA server.")

	utils.RegisterRecovery(app)
	utils.RegisterLogger(app)
	utils.RegisterViews(app)
	utils.RegisterErrorHandlers(app)

	server := &Server{
		App: app,
		Handlers: handlers.New(services),
	}

	server.setupRoutes()

	return server
}

func (s Server) Start(addr string) {
	s.App.Logger().Info("Starting GEMA server.")

	s.App.Run(iris.Addr(addr), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:     true,
		DisablePathCorrection: true,
		DisableVersionChecker: true,
	}))
}

func (s *Server) setupRoutes() {
	s.App.Logger().Info("Setting up routes.")

	s.App.Get("/", s.Handlers.Dashboard.HQ)

	s.App.Get("/health", s.Handlers.Health)
	s.App.Post("/login", s.Handlers.LoginPost)
	s.App.Get("/setup", s.Handlers.SetupGet)
	s.App.Post("/setup", s.Handlers.SetupPost)

	dashRoute := s.App.Party("/dash")
	dashRoute.Get("/view", s.Handlers.Dashboard.HQ)
}

