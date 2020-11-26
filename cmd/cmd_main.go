package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/avrebarra/postaco/docbuilder/postman"
	"github.com/avrebarra/postaco/postaco"
	"github.com/avrebarra/postaco/server"
	"github.com/facebookgo/grace/gracehttp"
	"gopkg.in/go-playground/validator.v9"
)

type ConfigCommandMain struct {
	Port       int    `validate:"required,ne=0"`
	SourcePath string `validate:"required,ne="`
}

type CommandMain struct {
	config ConfigCommandMain
}

func NewCommandMain(cfg ConfigCommandMain) CommandMain {
	if err := validator.New().Struct(cfg); err != nil {
		panic(err)
	}
	return CommandMain{config: cfg}
}

func (c CommandMain) Run() error {
	pstc := postaco.New(postaco.Config{})
	err := pstc.BuildDir(context.Background(), postaco.ConfigBuildDir{
		SourcePath: c.config.SourcePath,
		OutputPath: filepath.Join(c.config.SourcePath, "postaco"),
		DocTitle:   "Postaco - Documentation Server",
		Logger:     os.Stdout,

		PostmanDocBuilder: postman.New(postman.Config{}),
	})
	if err != nil {
		return err
	}

	svr := server.NewServer(server.ConfigServer{
		Path: filepath.Join(c.config.SourcePath, "postaco"),
	})
	defer svr.Close()

	// start server
	finalport := fmt.Sprintf(":%d", c.config.Port)
	fmt.Printf("documentation server started in http://localhost%s", finalport)
	err = gracehttp.Serve(&http.Server{Addr: finalport, Handler: svr.GetHandler()})
	if err != nil {
		return nil
	}

	return nil
}
