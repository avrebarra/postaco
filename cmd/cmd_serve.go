package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/avrebarra/postaco/server"
	"github.com/facebookgo/grace/gracehttp"
	"gopkg.in/go-playground/validator.v9"
)

type ConfigCommandServe struct {
	Port int    `validate:"required,ne=0"`
	Mode string `validate:"required,ne="`
}

type CommandServe struct {
	config ConfigCommandServe
}

func NewCommandServe(cfg ConfigCommandServe) CommandServe {
	if err := validator.New().Struct(cfg); err != nil {
		panic(err)
	}
	return CommandServe{config: cfg}
}

func (c CommandServe) Run() error {
	svr := server.NewServer(server.ConfigServer{
		Mode: c.config.Mode,
	})
	defer svr.Close()

	// start server
	finalport := fmt.Sprintf(":%d", c.config.Port)
	fmt.Printf("starting server in http://localhost:%s", finalport)
	err := gracehttp.Serve(&http.Server{Addr: finalport, Handler: svr.Router})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
