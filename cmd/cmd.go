package cmd

import (
	"path"

	"github.com/avrebarra/postaco/docbuilder/postman"
	"github.com/leaanthony/clir"
)

var cmd *clir.Cli

func customBanner(cli *clir.Cli) string {
	return ``
}

func Initialize() {
	cmd = clir.NewCli("postaco", "service", "v1")
	// Command.SetBannerFunction(customBanner)

	quiet := false
	cmd.BoolFlag("quiet", "perform quiet operation", &quiet)

	// default cli action on no params / flag
	cmd.Action(func() (err error) {
		cmd.PrintHelp()
		cmd := NewCommandBuild(ConfigCommandBuild{
			Quiet:      quiet,
			SourcePath: "tmp/collections",
			OutputPath: "tmp/.devweb",

			PostmanDocBuilder: postman.New(postman.Config{}),
		})
		return cmd.Run()
	})

	cmdBuild := cmd.NewSubCommand("build", "build documentation folder")
	{
		src := "."
		out := "."
		cmdBuild.StringFlag("src", "source folder (default: current dir)", &src)
		cmdBuild.StringFlag("out", "output folder (default: current dir)", &out)
		cmdBuild.Action(func() (err error) {
			srcclean := path.Clean(src)
			outclean := path.Clean(out)

			cmd := NewCommandBuild(ConfigCommandBuild{
				Quiet:      quiet,
				SourcePath: srcclean,
				OutputPath: outclean,

				PostmanDocBuilder: postman.New(postman.Config{}),
			})
			return cmd.Run()
		})
	}

	cmdServer := cmd.NewSubCommand("server", "build and start documentation server")
	{
		port := 8877
		mode := "production"
		src := "."
		out := "."
		cmdServer.StringFlag("src", "source folder (default: current dir)", &src)
		cmdServer.StringFlag("out", "output folder (default: current dir)", &out)
		cmdServer.StringFlag("env", "environment to use [production|development] (default: production)", &mode)
		cmdServer.IntFlag("port", "port to bind (default: 8877)", &port)
		cmdServer.Action(func() (err error) {
			cmd := NewCommandServe(ConfigCommandServe{
				Port: port,
				Mode: mode,
			})
			return cmd.Run()
		})
	}
}

func Run() error {
	return cmd.Run()
}
