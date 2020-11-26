package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/avrebarra/postaco/docbuilder/postman"
	"github.com/avrebarra/postaco/postaco"
	"gopkg.in/go-playground/validator.v9"
)

type ConfigCommandBuild struct {
	Quiet      bool
	SourcePath string `validate:"required,ne="`
	OutputPath string `validate:"required,ne="`
}

type CommandBuild struct {
	config ConfigCommandBuild
}

func NewCommandBuild(cfg ConfigCommandBuild) CommandBuild {
	if err := validator.New().Struct(cfg); err != nil {
		panic(err)
	}
	return CommandBuild{config: cfg}
}

func (c CommandBuild) Log(msg string) {
	if !c.config.Quiet {
		fmt.Println(msg)
	}
}

func (c CommandBuild) Run() (err error) {
	pstc := postaco.New(postaco.Config{})
	err = pstc.BuildDir(context.Background(), postaco.ConfigBuildDir{
		SourcePath:        c.config.SourcePath,
		OutputPath:        c.config.OutputPath,
		DocTitle:          "Postaco - Documentation Server",
		Logger:            os.Stdout,
		PostmanDocBuilder: postman.New(postman.Config{}),
	})

	return nil
}
