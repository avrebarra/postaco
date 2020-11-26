package cmd

import (
	"fmt"
	"path"

	"github.com/leaanthony/clir"
)

var cmd *clir.Cli

func customBanner(cli *clir.Cli) string {
	return ``
}

func Initialize() {
	cmd = clir.NewCli("postaco", "service", "v1")
	// Command.SetBannerFunction(customBanner)

	port := 8877
	quiet := false
	cmd.BoolFlag("quiet", "perform quiet operation", &quiet)
	cmd.IntFlag("port", "port to bind (default: 8877)", &port)

	// default cli action on no params / flag
	cmd.Action(func() (err error) {
		if len(cmd.OtherArgs()) == 0 {
			err = fmt.Errorf("no path supplied")
			return
		}

		subcmd := NewCommandMain(ConfigCommandMain{
			Port:       port,
			SourcePath: cmd.OtherArgs()[0],
		})
		return subcmd.Run()
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

			subcmd := NewCommandBuild(ConfigCommandBuild{
				Quiet:      quiet,
				SourcePath: srcclean,
				OutputPath: outclean,
			})
			return subcmd.Run()
		})
	}
}

func Run() error {
	return cmd.Run()
}
