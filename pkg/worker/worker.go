package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/cron"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/telegram/commands"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Worker struct {
	bot     *telebot.Bot
	jobs    *cron.Jobs
	db      database.Database
	youtube ytube.Service
	cfg     config.AppConfig
	stop    chan error
}

func (w *Worker) Start(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)

	defer cancel()

	w.stop = make(chan error, 1)

	defer close(w.stop)

	go func() {
		for {
			select {
			case <-ctx.Done():
				zerolog.Ctx(ctx).Warn().Err(ctx.Err()).Msg("Context done")
				cancel()

				return
			case err := <-w.stop:
				if err != nil {
					zerolog.Ctx(ctx).Warn().Err(err).Msg("Stop signal")
					cancel()
				}
			}
		}
	}()

	var wg sync.WaitGroup

	wg.Add(3) //nolint:gomnd

	go w.JobsWorker(ctx, &wg)
	go w.DatabaseWorker(ctx, &wg)
	go w.TelegramWorker(ctx, &wg)

	wg.Wait()

	if err := w.db.DB().Close(); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Fail to close database")
	}

	zerolog.Ctx(ctx).Warn().Msg("Worker stopped")
}

func (w *Worker) JobsWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if w.jobs == nil {
		zerolog.Ctx(ctx).Warn().Msg("Cron jobs are disabled")

		return
	}

	err := w.jobs.Start(ctx)

	switch err { //nolint:errorlint
	case context.Canceled:
	case nil:
		return
	default:
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Fail start jobs")
		w.stop <- err

		return
	}
}

func (w *Worker) DatabaseWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	zerolog.Ctx(ctx).Info().Msg("Starting database worker...")

	w.db.Worker(ctx)

	zerolog.Ctx(ctx).Warn().Msg("Database worker stopped")
}

func (w *Worker) TelegramWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := zerolog.Ctx(ctx)

	w.bot.OnError = func(err error, tx telebot.Context) {
		_ = tx.Reply(fmt.Sprintf("Error: %s", err.Error()))

		logger.Error().Err(err).Msg("Bot error")
	}

	go func() {
		logger.Info().Msg("Starting telegram bot...")

		cmds := commands.New(w.cfg, w.youtube, w.db)

		err := cmds.Setup(ctx, w.bot)
		if err != nil {
			logger.Warn().Err(err).Msg("Fail start telegram bot")
			w.stop <- err

			return
		}

		w.bot.Start()
	}()

	<-ctx.Done()
	logger.Warn().Err(ctx.Err()).Msg("Stopping telegram bot...")
	w.bot.Stop()
	logger.Warn().Msg("Telegram stopped")
}
