package worker

import (
	"context"
	"fmt"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/cron"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var Worker = &cli.Command{
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
		cfg := config.Ctx(cmd.Context)
		logger := zerolog.Ctx(cmd.Context).
			With().
			Str("context", "worker").
			Logger()

		logger.Debug().Msg("Creating youtube service instance")

		youtube, err := ytube.NewService(cmd.Context, cfg.Youtube)
		if err != nil {
			logger.Error().Err(err).Msg("Starting error")

			return err
		}

		logger.Debug().Msg("Creating bot instance")

		bot, err := telegram.NewBot(cmd.Context, *cfg, youtube)
		if err != nil {
			logger.Error().Err(err).Msg("Starting error")

			return err
		}

		ctx, cancel := support.WithKillSignal(cmd.Context)

		go func() {
			<-ctx.Done()
			logger.Warn().Msg("Stoping bot...")
			bot.Stop()
			logger.Warn().Msg("Stoped...")
		}()

		logger.Info().Msg("Starting bot...")

		bot.OnError = func(err error, tctx telebot.Context) {
			_ = tctx.Reply(fmt.Sprintf("Error: %s", err.Error()))

			logger.Error().Err(err).Msg("Bot error")
		}

		go func() {
			if !cmd.Bool("cron") {
				logger.Debug().Msg("cron disabled")

				return
			}

			logger.Info().Msg("cron enabled")

			jobs, err := cron.New(ctx, *cfg, bot, youtube)
			if err != nil {
				logger.Warn().Err(err).Msg("Fail to create cron jobs")
				cancel()
			}

			err = jobs.Start(ctx)

			switch err { //nolint:errorlint
			case context.Canceled:
			case nil:
				return
			default:
				logger.Warn().Err(err).Msg("Fail start jobs")
				cancel()
			}
		}()

		bot.Start()

		logger.Info().Msg("Done")

		return nil
	},
}
