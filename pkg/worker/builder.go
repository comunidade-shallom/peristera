package worker

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/cron"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
)

func New(ctx context.Context, cfg config.AppConfig, enableCron bool) (*Worker, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("fn", "worker:builder").
		Logger()

	serv := &Worker{}

	var err error

	serv.db, err = database.Open(cfg.Store.Path)
	if err != nil {
		return serv, err
	}

	logger.Info().Msg("Database opened")

	youtube, err := ytube.NewService(ctx, cfg.Youtube)
	if err != nil {
		return serv, err
	}

	logger.Info().Msg("YouTube service created")

	serv.bot, err = telegram.NewBot(cfg)
	if err != nil {
		return serv, err
	}

	logger.Info().Msg("Telegram bot created")

	if enableCron {
		serv.jobs, err = cron.New(ctx, cron.Options{
			Bot:      serv.bot,
			Database: serv.db,
			YouTube:  youtube,
			Config:   cfg,
		})
		if err != nil {
			return serv, err
		}

		logger.Info().Msg("Jobs service created")
	} else {
		logger.Warn().Msg("Jobs service disabled")
	}

	return serv, nil
}
