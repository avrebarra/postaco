package server

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type ConfigServer struct {
	Path string `validate:"required"`
}

type Server struct {
	config       ConfigServer
	router       http.Handler
	closerJaeger io.Closer
}

// NewServer creates new server instance
func NewServer(cfg ConfigServer) (s Server) {
	if err := validator.New().Struct(cfg); err != nil {
		panic(err)
	}

	s = Server{
		config: cfg,
		router: echo.New(),
	}

	// setup router
	s.SetupRouter()

	return
}

func (s *Server) SetupRouter() {
	router := s.router.(*echo.Echo)
	router.Validator = &CustomValidator{Validator: validator.New()}

	router.Static("/", s.config.Path)

	// Setup middlewares and add-ons
	router.Use(middleware.RemoveTrailingSlash())
	router.Use(middleware.RequestID())
	router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
}

func (s *Server) GetHandler() http.Handler {
	return s.router
}

func (s *Server) Close() error {
	return nil
}
