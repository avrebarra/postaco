package server

import (
	"io"
	"log"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type ConfigServer struct {
	Mode string
}

type Server struct {
	ConfigServer `validate:"required,structonly"`
	Router       http.Handler `validate:"required,structonly"`
	closerJaeger io.Closer
}

// NewServer creates new server instance
func NewServer(cfg ConfigServer) (s Server) {
	s = Server{
		ConfigServer: cfg,
		Router:       echo.New(),
	}

	// assert dependencies
	if err := validator.New().Struct(s.ConfigServer); err != nil {
		panic(err)
	}
	if err := validator.New().Struct(s); err != nil {
		log.Fatal(err)
	}

	// setup router
	s.SetupRouter()

	return
}

func (s *Server) SetupRouter() {
	router := s.Router.(*echo.Echo)
	router.Validator = &CustomValidator{Validator: validator.New()}
	router.Static("/", ".")

	// - development endpoints
	if s.ConfigServer.Mode == "development" {
		router.Any("/debug/statsviz/ws", echo.WrapHandler(http.HandlerFunc(statsviz.Ws)))
		router.Any("/debug/statsviz//ws", echo.WrapHandler(http.HandlerFunc(statsviz.Ws)))
		router.Any("/debug/statsviz/*", echo.WrapHandler(statsviz.Index))
	}

	// Setup middlewares and add-ons
	router.Use(middleware.RemoveTrailingSlash())
	router.Use(middleware.RequestID())
	if s.ConfigServer.Mode == "production" {
		router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
	}

	router.HTTPErrorHandler = ErrorHandler
}

func (s *Server) GetHandler() http.Handler {
	return s.Router
}

func (s *Server) Close() error {
	return nil
}
