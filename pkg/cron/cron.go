package cron

import (
	"context"
	"sync"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/telegram/sender"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type MessageBuilder func(sender.Chats) ([]sender.Sendable, error)

type Jobs struct {
	cfg       config.AppConfig
	youtube   ytube.Service
	bot       *telebot.Bot
	scheduler *gocron.Scheduler
	senderCh  chan sender.Sendable
}

var ErrNoChatToBroadcast = errors.Business("No chats to broadcast", "CRON:001")

func (j *Jobs) Start(ctx context.Context) error {
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		err, _ := recoverData.(error)

		zerolog.Ctx(ctx).
			Error().
			Str("context", "jobs:panic-handler").
			Str("job-name", jobName).
			Err(err).
			Stack().
			Msg("Panic to run job")
	})

	j.senderCh = make(chan sender.Sendable, 2) //nolint:gomnd

	var wg sync.WaitGroup

	wg.Add(1) // SendableWorker

	go func() {
		sender.SendableWorker(ctx, j.senderCh, j.bot)
		wg.Done()
	}()

	j.scheduler.StartAsync()

	<-ctx.Done()

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "jobs:manager").
		Logger()

	logger.Warn().Err(ctx.Err()).Msg("Stopping cron jobs...")

	j.scheduler.Stop()

	close(j.senderCh)

	logger.Warn().Msg("Cron jobs are stopped")

	wg.Wait()

	logger.Warn().Msg("Waiting for workers...")

	return nil
}

func (j *Jobs) register(ctx context.Context) error {
	timezone, err := time.LoadLocation(j.cfg.Timezone)
	if err != nil {
		return err
	}

	j.scheduler = gocron.NewScheduler(timezone)

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "jobs:register").
		Logger()

	if len(j.cfg.Cron.LastUpdates) == 0 {
		logger.Warn().Msg("No cron entries from last updates")

		return nil
	}

	logger.Info().Msg("Registring last updates")

	lastUpdates := func() error {
		return j.LastVideos(ctx)
	}

	j.scheduler.SingletonModeAll()

	for _, entry := range j.cfg.Cron.LastUpdates {
		_, err := j.scheduler.
			Cron(entry).
			Do(lastUpdates)
		if err != nil {
			return err
		}

		logger.Info().Msgf("Registred: %s", entry)
	}

	return nil
}

func (j *Jobs) jobLogger(ctx context.Context, job string) zerolog.Logger {
	return zerolog.Ctx(ctx).
		With().
		Str("context", "jobs").
		Str("job", job).
		Logger()
}

func (j *Jobs) broadcast(fn MessageBuilder, opts ...interface{}) error {
	if len(j.cfg.Telegram.Broadcast) == 0 {
		return ErrNoChatToBroadcast
	}

	entries, err := fn(j.cfg.Telegram.Broadcast)
	if err != nil {
		return err
	}

	for _, msg := range entries {
		j.senderCh <- msg
	}

	return nil
}
