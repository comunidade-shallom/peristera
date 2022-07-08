package main

import (
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/worker"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var WorkerCmd = &cli.Command{
	Name:  "worker",
	Usage: "Start telegram bot worker",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "cron",
			Usage:       "enable cron service",
			DefaultText: "false",
		},
	},
	Action: func(cmd *cli.Context) error {
		ctx, cancel := support.WithKillSignal(cmd.Context)

		defer cancel()

		cfg := *config.Ctx(ctx)

		wk, err := worker.New(ctx, cfg, cmd.Bool("cron"))
		if err != nil {
			return err
		}

		wk.Start(ctx)

		zerolog.Ctx(ctx).Info().Err(ctx.Err()).Msg("Worker done")

		return nil
	},
}
