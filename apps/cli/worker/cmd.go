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
	Action: func(c *cli.Context) error {
		cfg := config.Ctx(c.Context)
		logger := zerolog.Ctx(c.Context).
			With().
			Str("context", "worker").
			Logger()

		logger.Debug().Msg("Creating youtube service instance")

		youtube, err := ytube.NewService(c.Context, cfg.Youtube)

		if err != nil {
			logger.Error().Err(err).Msg("Starting error")

			return err
		}

		logger.Debug().Msg("Creating bot instance")
		bot, err := telegram.NewBot(c.Context, *cfg, youtube)

		if err != nil {
			logger.Error().Err(err).Msg("Starting error")

			return err
		}

		ctx, cancel := support.WithKillSignal(c.Context)

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
			jobs, err := cron.New(ctx, *cfg, bot, youtube)
			if err != nil {
				logger.Warn().Err(err).Msg("Fail to create cron jobs")
				cancel()
			}

			switch jobs.Start(ctx) {
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
