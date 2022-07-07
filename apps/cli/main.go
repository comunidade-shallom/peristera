package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/comunidade-shallom/peristera/apps/cli/system"
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/support/logger"
	"github.com/pterm/pterm"
	zero "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Description:          "Peristera - Telegram BOT",
		Usage:                "Peristera CLI",
		Version:              config.Version(),
		Copyright:            "https://github.com/comunidade-shallom",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from",
				DefaultText: fmt.Sprintf("%s/peristera.yml", support.GetBinDirPath()),
			},
			&cli.BoolFlag{
				Name:        "no-banner",
				Usage:       "hide initial banner",
				DefaultText: "false",
			},
			&cli.StringFlag{
				Name:        "level",
				Aliases:     []string{"l"},
				Usage:       "define log level",
				DefaultText: "info",
			},
		},
		Commands: []*cli.Command{WorkerCmd, system.SystemCmd},
		Before:   beforeRun,
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("Peristera - Telegram BOT")
		fmt.Println("")
		fmt.Println(config.VersionVerbose())
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		zero.Fatal().Err(err).Msg("Fail run application")
		os.Exit(1)
	}
}

func beforeRun(ctx *cli.Context) error {
	pterm.Debug.Debugger = !ctx.Bool("debug")

	if !ctx.Bool("no-banner") {
		pterm.DefaultHeader.
			WithMargin(5). //nolint:gomnd
			Println("Peristera CLI \n" + config.Version())
	}

	appConfig, err := config.Load(ctx.String("config"))

	logLevel := ctx.String("level")

	if appConfig.Logger.Debug(logLevel) {
		appConfig.Debug = true
	}

	if err != nil {
		e, ok := err.(errors.BusinessError) //nolint:errorlint
		if ok && e.ErrorCode == config.ConfigFileWasCreated.ErrorCode {
			zero.Warn().Msg(err.Error())
		} else {
			zero.Fatal().Err(err).Msg("Fail to load config")

			return err
		}
	}

	ctx.Context = appConfig.WithContext(ctx.Context)

	logger.SetupLogger(appConfig, logLevel)

	log := logger.Logger("", appConfig.Tags())

	ctx.Context = log.WithContext(appConfig.WithContext(ctx.Context))

	return nil
}
