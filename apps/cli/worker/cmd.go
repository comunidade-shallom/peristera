package worker

import (
	"context"
	"fmt"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/cron"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/comunidade-shallom/peristera/pkg/telegram/commands"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

type serviceContainer struct {
	cfg     config.AppConfig
	bot     *telebot.Bot
	youtube ytube.Service
	logger  zerolog.Logger
}

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
		ctx, cancel := support.WithKillSignal(cmd.Context)

		defer cancel()

		cfg := *config.Ctx(cmd.Context)
		logger := zerolog.Ctx(cmd.Context).
			With().
			Str("context", "worker").
			Logger()

		var err error

		services := serviceContainer{
			cfg:    cfg,
			logger: logger,
		}

		logger.Debug().Msg("Creating youtube service instance")

		services.youtube, err = ytube.NewService(cmd.Context, cfg.Youtube)
		if err != nil {
			return err
		}

		services.bot, err = initBot(ctx, services)
		if err != nil {
			return err
		}

		if cmd.Bool("cron") {
			go func() {
				errSig := initCron(ctx, services)
				err := <-errSig
				if err != nil {
					logger.Error().Err(err).Msg("error signal from cron")
				}
			}()
		} else {
			logger.Warn().Msg("Cron is disabled")
		}

		<-ctx.Done()

		logger.Info().Err(ctx.Err()).Msg("Done")

		return nil
	},
}

func initBot(ctx context.Context, services serviceContainer) (*telebot.Bot, error) {
	logger := services.logger

	logger.Debug().Msg("Creating bot instance")

	bot, err := telegram.NewBot(services.cfg)
	if err != nil {
		return nil, err
	}

	err = commands.
		New(services.cfg, services.youtube).
		Setup(ctx, bot)

	if err != nil {
		return nil, err
	}

	bot.OnError = func(err error, tx telebot.Context) {
		_ = tx.Reply(fmt.Sprintf("Error: %s", err.Error()))

		logger.Error().Err(err).Msg("Bot error")
	}

	go func() {
		<-ctx.Done()
		logger.Warn().Err(ctx.Err()).Msg("Stoping bot...")
		bot.Stop()
		logger.Warn().Msg("Stoped...")
	}()

	go func() {
		logger.Info().Msg("Starting bot...")

		bot.Start()
	}()

	return bot, nil
}

func initCron(ctx context.Context, services serviceContainer) <-chan error {
	logger := services.logger

	out := make(chan error, 2) //nolint:gomnd

	logger.Info().Msg("Cron is enabled")
	logger.Debug().Msg("Creating cron jobs")

	jobs, err := cron.New(ctx, services.cfg, services.bot, services.youtube)
	if err != nil {
		logger.Warn().Err(err).Msg("Fail to create cron jobs")
		out <- err
		close(out)

		return out
	}

	go func() {
		defer close(out)

		err = jobs.Start(ctx)

		switch err { //nolint:errorlint
		case context.Canceled:
		case nil:
			return
		default:
			logger.Warn().Err(err).Msg("Fail start jobs")
			out <- err

			return
		}
	}()

	return out
}
